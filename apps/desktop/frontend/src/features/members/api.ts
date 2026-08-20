import { desktopApp } from "../../platform/wails";

export type Member = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  identificationNumber: string;
  dateOfBirth: string;
  address: string;
  status: string;
};

export type MemberInput = Omit<Member, "id">;

type MemberBindings = {
  ListMembers(): Promise<Member[]>;
  CreateMember(input: MemberInput): Promise<Member>;
  UpdateMember(id: string, input: MemberInput): Promise<Member>;
  ArchiveMember(id: string): Promise<void>;
};

function members(): MemberBindings {
  const bridge = desktopApp() as MemberBindings | null;
  if (!bridge) {
    throw new Error("Member management is only available in the desktop app");
  }
  return bridge;
}

export function listMembers(): Promise<Member[]> {
  return members().ListMembers();
}

export function createMember(input: MemberInput): Promise<Member> {
  return members().CreateMember(input);
}

export function updateMember(id: string, input: MemberInput): Promise<Member> {
  return members().UpdateMember(id, input);
}

export function archiveMember(id: string): Promise<void> {
  return members().ArchiveMember(id);
}
