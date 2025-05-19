export type LoginPayload = {
  email: string;
  password: string;
};

export type SignupPayload = {
  email: string;
  password: string;
};

export type AuthResponse = {
  access_token: string;
  refresh_token: string;
};
