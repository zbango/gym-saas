import { useEffect, useState } from "react";
import { isDesktopApp } from "../../platform/wails";
import { archiveMember, createMember, listMembers, type Member, type MemberInput, updateMember } from "./api";

const emptyMemberInput: MemberInput = {
  firstName: "",
  lastName: "",
  email: "",
  phone: "",
  identificationNumber: "",
  dateOfBirth: "",
  address: "",
  status: "active"
};

export function useMembers() {
  const [members, setMembers] = useState<Member[]>([]);
  const [form, setForm] = useState<MemberInput>(emptyMemberInput);
  const [editingMemberID, setEditingMemberID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const desktopRuntime = isDesktopApp();

  async function refresh() {
    if (!desktopRuntime) {
      return;
    }
    try {
      setMembers(await listMembers());
      setError(null);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to load members");
    }
  }

  useEffect(() => {
    void refresh();
  }, []);

  function setField(field: keyof MemberInput, value: string) {
    setForm((current) => ({ ...current, [field]: value }));
  }

  function startEditing(member: Member) {
    setEditingMemberID(member.id);
    setForm({
      firstName: member.firstName,
      lastName: member.lastName,
      email: member.email,
      phone: member.phone,
      identificationNumber: member.identificationNumber,
      dateOfBirth: member.dateOfBirth,
      address: member.address,
      status: member.status
    });
    setError(null);
  }

  function cancelEditing() {
    setEditingMemberID(null);
    setForm(emptyMemberInput);
    setError(null);
  }

  async function save() {
    if (!desktopRuntime) {
      setError("Member management is only available in the desktop app");
      return false;
    }
    setSaving(true);
    setError(null);
    try {
      if (editingMemberID) {
        await updateMember(editingMemberID, form);
      } else {
        await createMember(form);
      }
      cancelEditing();
      await refresh();
      return true;
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to save member");
      return false;
    } finally {
      setSaving(false);
    }
  }

  async function archive(member: Member) {
    setError(null);
    try {
      await archiveMember(member.id);
      if (editingMemberID === member.id) {
        cancelEditing();
      }
      await refresh();
      return true;
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to archive member");
      return false;
    }
  }

  return { members, form, editingMemberID, error, saving, setField, startEditing, cancelEditing, save, archive };
}
