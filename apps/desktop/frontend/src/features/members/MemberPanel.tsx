import { useState } from "react";
import { Panel, ShellButton } from "@gym-saas/ui";
import type { Member } from "./api";
import { useMembers } from "./useMembers";

const memberInputStyle = {
  padding: "10px 12px",
  background: "var(--gs-input-background)",
  color: "var(--gs-input-text)",
  border: "1px solid var(--gs-border)",
  borderRadius: 12,
  font: "inherit"
};

export function MemberPanel() {
  const { members, form, editingMemberID, error, saving, setField, startEditing, cancelEditing, save, archive } = useMembers();
  const [archiveCandidate, setArchiveCandidate] = useState<Member | null>(null);
  const [archiving, setArchiving] = useState(false);

  return (
    <Panel title="Members" eyebrow="Member directory">
      <p style={{ marginTop: 0, color: "var(--gs-text-muted)" }}>
        {editingMemberID ? "Update this member's profile." : "Add a member to your local gym."}
      </p>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(210px, 1fr))", gap: 12 }}>
        <MemberField label="First name" value={form.firstName} onChange={(value) => setField("firstName", value)} />
        <MemberField label="Last name" value={form.lastName} onChange={(value) => setField("lastName", value)} />
        <MemberField label="Phone" value={form.phone} onChange={(value) => setField("phone", value)} />
        <MemberField label="Email" type="email" value={form.email} onChange={(value) => setField("email", value)} />
        <MemberField label="Date of birth" type="date" value={form.dateOfBirth} onChange={(value) => setField("dateOfBirth", value)} />
        <MemberField label="Address" value={form.address} onChange={(value) => setField("address", value)} />
        <MemberField label="Identification number" value={form.identificationNumber} onChange={(value) => setField("identificationNumber", value)} />
        <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
          Status
          <select value={form.status} onChange={(event) => setField("status", event.target.value)} style={memberInputStyle}>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="blocked">Blocked</option>
          </select>
        </label>
      </div>
      {error ? <p style={{ color: "var(--gs-danger, #b42318)", marginBottom: 0 }}>{error}</p> : null}
      <div style={{ display: "flex", gap: 8, marginTop: 16, flexWrap: "wrap" }}>
        <ShellButton onClick={() => void save()} disabled={saving}>
          {saving ? "Saving..." : editingMemberID ? "Save changes" : "Create member"}
        </ShellButton>
        {editingMemberID ? <ShellButton variant="secondary" onClick={cancelEditing}>Cancel</ShellButton> : null}
      </div>
      <div style={{ display: "grid", gap: 8, marginTop: 20 }}>
        {members.length === 0 ? (
          <p style={{ color: "var(--gs-text-muted)", margin: 0 }}>No active members yet.</p>
        ) : members.map((member) => (
          <div key={member.id} style={{ display: "flex", gap: 12, justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", borderTop: "1px solid var(--gs-border)", paddingTop: 12 }}>
            <div>
              <strong>{member.firstName} {member.lastName}</strong>
              <div style={{ color: "var(--gs-text-muted)", fontSize: 14 }}>{member.phone}{member.email ? ` · ${member.email}` : ""} · {member.status}</div>
            </div>
            <div style={{ display: "flex", gap: 8 }}>
              <ShellButton variant="secondary" onClick={() => startEditing(member)}>Edit</ShellButton>
              <ShellButton variant="secondary" onClick={() => setArchiveCandidate(member)}>Archive</ShellButton>
            </div>
          </div>
        ))}
      </div>
      {archiveCandidate ? (
        <div style={{ marginTop: 20, padding: 16, borderRadius: 12, background: "var(--gs-accent-soft)", display: "grid", gap: 10 }}>
          <strong>Archive {archiveCandidate.firstName} {archiveCandidate.lastName}?</strong>
          <span>This removes them from the active list but does not permanently delete their record.</span>
          <div style={{ display: "flex", gap: 8 }}>
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
    <label style={{ display: "grid", gap: 6, fontWeight: 600 }}>
      {props.label}
      <input type={props.type ?? "text"} value={props.value} onChange={(event) => props.onChange(event.target.value)} style={memberInputStyle} />
    </label>
  );
}
