export type UserRole = "super_admin" | "gym_admin" | "receptionist";

export type AuthenticatedUser = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
};

export type LoginCredentials = {
  username: string;
  password: string;
};

type MockAccount = AuthenticatedUser & {
  username: string;
  password: string;
};

const mockAccounts: readonly MockAccount[] = [
  {
    id: "mock-super-admin",
    name: "Sofía Ramírez",
    email: "superadmin@zeus.gym",
    role: "super_admin",
    username: "superadmin",
    password: "demo"
  },
  {
    id: "mock-gym-admin",
    name: "Mateo León",
    email: "admin@zeus.gym",
    role: "gym_admin",
    username: "admin",
    password: "demo"
  },
  {
    id: "mock-receptionist",
    name: "Carla Mora",
    email: "recepcion@zeus.gym",
    role: "receptionist",
    username: "recepcion",
    password: "demo"
  }
];

export const mockLoginAccounts = mockAccounts.map(({ username, password, role }) => ({ username, password, role }));

export async function authenticateMockUser(credentials: LoginCredentials): Promise<AuthenticatedUser | null> {
  const account = mockAccounts.find(
    ({ username, password }) =>
      username === credentials.username.trim().toLowerCase() && password === credentials.password
  );

  if (!account) {
    return null;
  }

  const { username: _username, password: _password, ...user } = account;
  return user;
}
