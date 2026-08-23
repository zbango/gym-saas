import { useState } from "react";
import { Panel, ShellButton } from "@gym-saas/ui";
import type { Member } from "./api";
import { useMembers } from "./useMembers";

const fieldClassName = "min-h-11 rounded-xl border border-[var(--color-brand-border)] bg-[var(--color-brand-input)] px-3 py-2.5 text-[var(--color-brand-text)] outline-0";
const labelClassName = "grid gap-1.5 font-semibold";

export function MemberPanel() {
  const { members, form, editingMemberID, error, saving, setField, startEditing, cancelEditing, save, archive } = useMembers();
  const [archiveCandidate, setArchiveCandidate] = useState<Member | null>(null);
  const [archiving, setArchiving] = useState(false);

  return (
    <Panel title="Members" eyebrow="Member directory">
      <p className="mt-0 text-[var(--color-brand-muted)]">
        {editingMemberID ? "Update this member's profile." : "Add a member to your local gym."}
      </p>
      <div className="grid gap-3 [grid-template-columns:repeat(auto-fit,minmax(210px,1fr))]">
        <MemberField label="First name" value={form.firstName} onChange={(value) => setField("firstName", value)} />
        <MemberField label="Last name" value={form.lastName} onChange={(value) => setField("lastName", value)} />
        <MemberField label="Phone" value={form.phone} onChange={(value) => setField("phone", value)} />
        <MemberField label="Email" type="email" value={form.email} onChange={(value) => setField("email", value)} />
        <MemberField label="Date of birth" type="date" value={form.dateOfBirth} onChange={(value) => setField("dateOfBirth", value)} />
        <MemberField label="Address" value={form.address} onChange={(value) => setField("address", value)} />
        <MemberField label="Identification number" value={form.identificationNumber} onChange={(value) => setField("identificationNumber", value)} />
        <label className={labelClassName}>
          Status
          <select className={fieldClassName} value={form.status} onChange={(event) => setField("status", event.target.value)}>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="blocked">Blocked</option>
          </select>
        </label>
      </div>
      {error ? <p className="mb-0 text-[var(--color-brand-danger)]">{error}</p> : null}
      <div className="mt-4 flex flex-wrap gap-2">
        <ShellButton onClick={() => void save()} disabled={saving}>
          {saving ? "Saving..." : editingMemberID ? "Save changes" : "Create member"}
        </ShellButton>
        {editingMemberID ? <ShellButton variant="secondary" onClick={cancelEditing}>Cancel</ShellButton> : null}
      </div>
      <div className="mt-5 grid gap-2">
        {members.length === 0 ? (
          <p className="m-0 text-[var(--color-brand-muted)]">No active members yet.</p>
        ) : members.map((member) => (
          <div key={member.id} className="flex flex-wrap items-center justify-between gap-3 border-t border-[var(--color-brand-border)] pt-3">
            <div>
              <strong>{member.firstName} {member.lastName}</strong>
              <div className="text-sm text-[var(--color-brand-muted)]">{member.phone}{member.email ? ` · ${member.email}` : ""} · {member.status}</div>
            </div>
            <div className="flex gap-2">
              <ShellButton variant="secondary" onClick={() => startEditing(member)}>Edit</ShellButton>
              <ShellButton variant="secondary" onClick={() => setArchiveCandidate(member)}>Archive</ShellButton>
            </div>
          </div>
        ))}
      </div>
      {archiveCandidate ? (
        <div className="mt-5 grid gap-2.5 rounded-xl bg-[var(--color-brand-accent-soft)] p-4">
          <strong>Archive {archiveCandidate.firstName} {archiveCandidate.lastName}?</strong>
          <span>This removes them from the active list but does not permanently delete their record.</span>
          <div className="flex gap-2">
            <ShellButton variant="secondary" onClick={() => setArchiveCandidate(null)} disabled={archiving}>Cancel</ShellButton>
            <ShellButton
              disabled={archiving}
              onClick={async () => {
                setArchiving(true);
                try {
                  if (await archive(archiveCandidate)) {
                    setArchiveCandidate(null);
                  }
                } finally {
                  setArchiving(false);
                }
              }}
            >
              {archiving ? "Archiving..." : "Archive member"}
            </ShellButton>
          </div>
        </div>
      ) : null}
    </Panel>
  );
}

function MemberField(props: { label: string; value: string; onChange: (value: string) => void; type?: string }) {
  return (
    <label className={labelClassName}>
      {props.label}
      <input className={fieldClassName} type={props.type ?? "text"} value={props.value} onChange={(event) => props.onChange(event.target.value)} />
    </label>
  );
}
