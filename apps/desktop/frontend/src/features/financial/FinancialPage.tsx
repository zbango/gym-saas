import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { CashIcon, CheckCircleIcon, PageHeader, PlusIcon, RefreshIcon } from "@gym-saas/ui";
import { listMembers, type Member } from "../members/api";
import { isDesktopApp } from "../../platform/wails";
import {
  createPendingExpense,
  createPendingPayment,
  listExpenses,
  listPaymentsForMember,
  postExpense,
  postPayment,
  recordExpense,
  recordPayment,
  refundPayment,
  voidExpense,
  voidPayment,
  type Expense,
  type ExpenseInput,
  type Payment,
  type PaymentInput
} from "./api";

type Tab = "payments" | "expenses";

const paymentMethods = ["cash", "credit_card", "debit_card", "bank_transfer", "check", "other"] as const;

const emptyExpense: ExpenseInput = { amountCents: 0, currency: "USD", paymentMethod: "cash", reference: "", notes: "" };

export function FinancialPage() {
  const desktopRuntime = isDesktopApp();
  const [tab, setTab] = useState<Tab>("payments");
  const [members, setMembers] = useState<Member[]>([]);
  const [selectedMemberID, setSelectedMemberID] = useState("");
  const [payments, setPayments] = useState<Payment[]>([]);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [amount, setAmount] = useState("");
  const [method, setMethod] = useState<PaymentInput["paymentMethod"]>("cash");
  const [reference, setReference] = useState("");
  const [notes, setNotes] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const selectedMember = useMemo(() => members.find((member) => member.id === selectedMemberID), [members, selectedMemberID]);

  const refresh = useCallback(async () => {
    if (!desktopRuntime) {
      return;
    }
    setPending(true);
    try {
      const [loadedMembers, loadedExpenses] = await Promise.all([listMembers(), listExpenses()]);
      setMembers(loadedMembers);
      setExpenses(loadedExpenses);
      setSelectedMemberID((current) => current || loadedMembers[0]?.id || "");
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    } finally {
      setPending(false);
    }
  }, [desktopRuntime]);

  const refreshPayments = useCallback(async (memberID: string) => {
    if (!desktopRuntime || !memberID) {
      setPayments([]);
      return;
    }
    try {
      setPayments(await listPaymentsForMember(memberID));
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    }
  }, [desktopRuntime]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  useEffect(() => {
    void refreshPayments(selectedMemberID);
  }, [refreshPayments, selectedMemberID]);

  async function submitPayment(makePending: boolean) {
    const amountCents = centsFromInput(amount);
    if (!selectedMemberID || amountCents === null) {
      setError("Selecciona un miembro e ingresa un monto válido mayor a cero.");
      return;
    }
    setPending(true);
    try {
      const input: PaymentInput = { memberId: selectedMemberID, membershipId: "", kind: "other", amountCents, currency: "USD", paymentMethod: method, reference, notes };
      if (makePending) {
        await createPendingPayment(input);
      } else {
        await recordPayment(input);
      }
      setAmount("");
      setReference("");
      setNotes("");
      await refreshPayments(selectedMemberID);
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    } finally {
      setPending(false);
    }
  }

  async function submitExpense(makePending: boolean) {
    const amountCents = centsFromInput(amount);
    if (amountCents === null) {
      setError("Ingresa un monto válido mayor a cero.");
      return;
    }
    setPending(true);
    try {
      const input: ExpenseInput = { ...emptyExpense, amountCents, paymentMethod: method, reference, notes };
      if (makePending) {
        await createPendingExpense(input);
      } else {
        await recordExpense(input);
      }
      setAmount("");
      setReference("");
      setNotes("");
      await refresh();
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    } finally {
      setPending(false);
    }
  }

  async function changePayment(payment: Payment, action: "post" | "void" | "refund") {
    setPending(true);
    try {
      if (action === "post") await postPayment(payment.id);
      if (action === "void") await voidPayment(payment.id);
      if (action === "refund") await refundPayment(payment.id);
      await refreshPayments(selectedMemberID);
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    } finally {
      setPending(false);
    }
  }

  async function changeExpense(item: Expense, action: "post" | "void") {
    setPending(true);
    try {
      if (action === "post") await postExpense(item.id);
      if (action === "void") await voidExpense(item.id);
      await refresh();
      setError(null);
    } catch (cause) {
      setError(messageFor(cause));
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="px-16 py-14 text-[#172131] max-[1260px]:px-[43px] max-[1260px]:py-[46px] max-[920px]:px-[22px] max-[920px]:py-8 max-[560px]:px-[14px] max-[560px]:py-6">
      <PageHeader
        icon={<CashIcon />}
        title="Pagos y Gastos"
        description="Registra cobros de miembros y gastos operativos desde la base local del gimnasio."
        actions={<button className="inline-flex min-h-[50px] items-center gap-2 rounded-xl border-0 bg-[#374357] px-5 text-base font-extrabold text-white" type="button" onClick={() => void refresh()} disabled={pending}><RefreshIcon className="h-5 w-5" /> Actualizar</button>}
      />

      {!desktopRuntime ? <Message text="Los pagos y gastos se habilitan al ejecutar la aplicación de escritorio." /> : null}
      {error ? <Message text={error} error /> : null}

      <div className="mb-7 flex gap-3" role="tablist" aria-label="Finanzas">
        <TabButton active={tab === "payments"} onClick={() => setTab("payments")}>Cobros</TabButton>
        <TabButton active={tab === "expenses"} onClick={() => setTab("expenses")}>Gastos</TabButton>
      </div>

      {tab === "payments" ? (
        <section className="grid gap-7 xl:grid-cols-[360px_minmax(0,1fr)]">
          <FormCard title="Registrar cobro" submitLabel="Registrar cobro" pending={pending} onSubmit={() => void submitPayment(false)} secondaryLabel="Dejar pendiente" onSecondary={() => void submitPayment(true)}>
            <Field label="Miembro"><select value={selectedMemberID} onChange={(event) => setSelectedMemberID(event.target.value)}>{members.length === 0 ? <option value="">No hay miembros</option> : members.map((member) => <option key={member.id} value={member.id}>{member.firstName} {member.lastName}</option>)}</select></Field>
            <AmountField amount={amount} onChange={setAmount} />
            <MethodField method={method} onChange={setMethod} />
            <TextField label="Referencia" value={reference} onChange={setReference} />
            <TextField label="Notas" value={notes} onChange={setNotes} />
          </FormCard>
          <ListCard title={selectedMember ? `Cobros de ${selectedMember.firstName} ${selectedMember.lastName}` : "Cobros"} empty="Selecciona un miembro para ver sus cobros.">
            {payments.map((payment) => <PaymentRow key={payment.id} payment={payment} pending={pending} onAction={changePayment} />)}
          </ListCard>
        </section>
      ) : (
        <section className="grid gap-7 xl:grid-cols-[360px_minmax(0,1fr)]">
          <FormCard title="Registrar gasto" submitLabel="Registrar gasto" pending={pending} onSubmit={() => void submitExpense(false)} secondaryLabel="Dejar pendiente" onSecondary={() => void submitExpense(true)}>
            <AmountField amount={amount} onChange={setAmount} />
            <MethodField method={method} onChange={setMethod} />
            <TextField label="Referencia" value={reference} onChange={setReference} />
            <TextField label="Notas" value={notes} onChange={setNotes} />
          </FormCard>
          <ListCard title="Gastos operativos" empty="Todavía no hay gastos registrados.">
            {expenses.map((item) => <ExpenseRow key={item.id} expense={item} pending={pending} onAction={changeExpense} />)}
          </ListCard>
        </section>
      )}
    </div>
  );
}

function TabButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: string }) {
  return <button className={`rounded-full px-5 py-2.5 text-base font-extrabold ${active ? "bg-[#374357] text-white" : "bg-white text-[#59636e] ring-1 ring-[#dfe3e8]"}`} type="button" onClick={onClick}>{children}</button>;
}

function FormCard({ title, children, submitLabel, secondaryLabel, pending, onSubmit, onSecondary }: { title: string; children: ReactNode; submitLabel: string; secondaryLabel: string; pending: boolean; onSubmit: () => void; onSecondary: () => void }) {
  return <form className="grid h-fit gap-4 rounded-[22px] border border-[#e0e4e9] bg-white p-6 shadow-[0_3px_5px_rgba(20,31,48,0.04)]" onSubmit={(event) => { event.preventDefault(); onSubmit(); }}><h2 className="text-2xl font-extrabold">{title}</h2>{children}<button className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl border-0 bg-[linear-gradient(100deg,var(--color-brand-button-start)_0%,var(--color-brand-button-middle)_42%,var(--color-brand-button-end)_100%)] px-5 font-extrabold text-[var(--color-brand-button-text)] disabled:opacity-60" disabled={pending} type="submit"><PlusIcon className="h-5 w-5" /> {submitLabel}</button><button className="min-h-11 rounded-xl border border-[#cfd5dc] bg-white px-4 font-bold text-[#4e5965] disabled:opacity-60" disabled={pending} type="button" onClick={onSecondary}>{secondaryLabel}</button></form>;
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <label className="grid gap-2 text-sm font-bold text-[#4d5865]"><span>{label}</span>{children}</label>;
}

function AmountField({ amount, onChange }: { amount: string; onChange: (value: string) => void }) {
  return <Field label="Monto (USD)"><input inputMode="decimal" placeholder="0.00" value={amount} onChange={(event) => onChange(event.target.value)} /></Field>;
}

function MethodField({ method, onChange }: { method: string; onChange: (value: PaymentInput["paymentMethod"]) => void }) {
  return <Field label="Método"><select value={method} onChange={(event) => onChange(event.target.value as PaymentInput["paymentMethod"])}>{paymentMethods.map((item) => <option key={item} value={item}>{methodLabel(item)}</option>)}</select></Field>;
}

function TextField({ label, value, onChange }: { label: string; value: string; onChange: (value: string) => void }) {
  return <Field label={label}><input value={value} onChange={(event) => onChange(event.target.value)} /></Field>;
}

function ListCard({ title, children, empty }: { title: string; children: ReactNode; empty: string }) {
  const rows = Array.isArray(children) ? children : [children];
  return <section className="overflow-hidden rounded-[22px] border border-[#e0e4e9] bg-white shadow-[0_3px_5px_rgba(20,31,48,0.04)]"><header className="border-b border-[#e7eaee] px-6 py-5"><h2 className="text-2xl font-extrabold">{title}</h2></header>{rows.filter(Boolean).length === 0 ? <p className="p-8 text-[#68727e]">{empty}</p> : <div className="divide-y divide-[#e9ecf0]">{children}</div>}</section>;
}

function PaymentRow({ payment, pending, onAction }: { payment: Payment; pending: boolean; onAction: (payment: Payment, action: "post" | "void" | "refund") => void }) {
  return <article className="flex flex-wrap items-center gap-4 px-6 py-5"><div className="min-w-[180px] flex-1"><strong>{formatMoney(payment.amountCents, payment.currency)}</strong><p className="mt-1 text-sm text-[#68727e]">{payment.reference || "Sin referencia"} · {payment.paidAt ? formatDate(payment.paidAt) : "Pendiente"}</p></div><Status status={payment.status} />{payment.status === "pending" ? <><Action label="Cobrar" disabled={pending} onClick={() => onAction(payment, "post")} /><Action label="Anular" disabled={pending} onClick={() => onAction(payment, "void")} /></> : null}{payment.status === "posted" ? <Action label="Reembolsar" disabled={pending} onClick={() => onAction(payment, "refund")} /> : null}</article>;
}

function ExpenseRow({ expense, pending, onAction }: { expense: Expense; pending: boolean; onAction: (expense: Expense, action: "post" | "void") => void }) {
  return <article className="flex flex-wrap items-center gap-4 px-6 py-5"><div className="min-w-[180px] flex-1"><strong>{formatMoney(expense.amountCents, expense.currency)}</strong><p className="mt-1 text-sm text-[#68727e]">{expense.reference || "Sin referencia"} · {expense.paidAt ? formatDate(expense.paidAt) : "Pendiente"}</p></div><Status status={expense.status} />{expense.status === "pending" ? <><Action label="Pagar" disabled={pending} onClick={() => onAction(expense, "post")} /><Action label="Anular" disabled={pending} onClick={() => onAction(expense, "void")} /></> : null}</article>;
}

function Status({ status }: { status: string }) {
  return <span className={`rounded-full px-3 py-1 text-sm font-extrabold ${status === "posted" ? "bg-[var(--color-brand-success-soft)] text-[#427a4a]" : status === "pending" ? "bg-[#fff4d5] text-[#936d1e]" : "bg-[#f4e9e8] text-[var(--color-brand-danger)]"}`}>{statusLabel(status)}</span>;
}

function Action({ label, disabled, onClick }: { label: string; disabled: boolean; onClick: () => void }) {
  return <button className="rounded-lg border border-[#cfd5dc] bg-white px-3 py-2 text-sm font-bold text-[#374357] disabled:opacity-50" disabled={disabled} type="button" onClick={onClick}>{label}</button>;
}

function Message({ text, error = false }: { text: string; error?: boolean }) {
  return <p className={`mb-6 rounded-xl border px-5 py-4 ${error ? "border-[#f0c9c5] bg-[#fff3f2] text-[var(--color-brand-danger)]" : "border-[#dce3ea] bg-white text-[#59636e]"}`}>{text}</p>;
}

function centsFromInput(value: string): number | null {
  const parsed = Number(value.replace(",", "."));
  if (!Number.isFinite(parsed) || parsed <= 0) return null;
  const cents = Math.round(parsed * 100);
  return cents > 0 && Math.abs(parsed * 100 - cents) < 0.000001 ? cents : null;
}

function formatMoney(cents: number, currency: string): string { return new Intl.NumberFormat("en-US", { style: "currency", currency }).format(cents / 100); }
function formatDate(value: string): string { return new Intl.DateTimeFormat("es-EC", { dateStyle: "medium" }).format(new Date(value)); }
function methodLabel(value: string): string { return ({ cash: "Efectivo", credit_card: "Tarjeta de crédito", debit_card: "Tarjeta de débito", bank_transfer: "Transferencia", check: "Cheque", other: "Otro" } as Record<string, string>)[value] ?? value; }
function statusLabel(value: string): string { return ({ pending: "Pendiente", posted: "Registrado", voided: "Anulado", refunded: "Reembolsado" } as Record<string, string>)[value] ?? value; }
function messageFor(cause: unknown): string { return cause instanceof Error ? cause.message : "No se pudo completar la operación."; }
