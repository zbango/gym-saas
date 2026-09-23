import {
  CreatePendingPayment,
  ListPaymentsForMember,
  PostPayment,
  RecordPayment,
  RefundPayment,
  VoidPayment
} from "../../../wailsjs/go/main/PaymentAPI";
import {
  CreatePendingExpense,
  ListExpenses,
  PostExpense,
  RecordExpense,
  VoidExpense
} from "../../../wailsjs/go/main/ExpenseAPI";
import type { main } from "../../../wailsjs/go/models";

export type Payment = main.Payment;
export type PaymentInput = main.PaymentInput;
export type Expense = main.Expense;
export type ExpenseInput = main.ExpenseInput;

export function listPaymentsForMember(memberID: string): Promise<Payment[]> {
  return ListPaymentsForMember(memberID);
}

export function recordPayment(input: PaymentInput): Promise<Payment> {
  return RecordPayment(input);
}

export function createPendingPayment(input: PaymentInput): Promise<Payment> {
  return CreatePendingPayment(input);
}

export function postPayment(id: string): Promise<Payment> {
  return PostPayment(id);
}

export function voidPayment(id: string): Promise<Payment> {
  return VoidPayment(id);
}

export function refundPayment(id: string): Promise<Payment> {
  return RefundPayment(id);
}

export function listExpenses(): Promise<Expense[]> {
  return ListExpenses();
}

export function recordExpense(input: ExpenseInput): Promise<Expense> {
  return RecordExpense(input);
}

export function createPendingExpense(input: ExpenseInput): Promise<Expense> {
  return CreatePendingExpense(input);
}

export function postExpense(id: string): Promise<Expense> {
  return PostExpense(id);
}

export function voidExpense(id: string): Promise<Expense> {
  return VoidExpense(id);
}
