import {
  ArchiveMembershipPlan,
  CreateMembershipPlan,
  ListMembershipPlans,
  UpdateMembershipPlan
} from "../../../wailsjs/go/main/MembershipPlanAPI";
import type { main } from "../../../wailsjs/go/models";

export type MembershipPlan = Omit<main.MembershipPlan, "validityKind" | "durationUnit" | "status"> & {
  validityKind: "time" | "visits";
  durationUnit: "days" | "weeks" | "months" | "years";
  status: "active" | "inactive";
};

export type MembershipPlanInput = Omit<MembershipPlan, "id">;

export function listMembershipPlans(): Promise<MembershipPlan[]> {
  return ListMembershipPlans() as Promise<MembershipPlan[]>;
}

export function createMembershipPlan(input: MembershipPlanInput): Promise<MembershipPlan> {
  return CreateMembershipPlan(input) as Promise<MembershipPlan>;
}

export function updateMembershipPlan(id: string, input: MembershipPlanInput): Promise<MembershipPlan> {
  return UpdateMembershipPlan(id, input) as Promise<MembershipPlan>;
}

export function archiveMembershipPlan(id: string): Promise<void> {
  return ArchiveMembershipPlan(id);
}
