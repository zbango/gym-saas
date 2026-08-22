import {
  ArchiveMember,
  CreateMember,
  ListMembers,
  UpdateMember
} from "../../../wailsjs/go/main/MemberAPI";
import type { main } from "../../../wailsjs/go/models";

export type Member = main.Member;
export type MemberInput = main.MemberInput;

export function listMembers(): Promise<Member[]> {
  return ListMembers();
}

export function createMember(input: MemberInput): Promise<Member> {
  return CreateMember(input);
}

export function updateMember(id: string, input: MemberInput): Promise<Member> {
  return UpdateMember(id, input);
}

export function archiveMember(id: string): Promise<void> {
  return ArchiveMember(id);
}
