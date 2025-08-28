export type AuthState = {
    token: string;
}

export type SignInRequest = {
    password: string;
}

export type LoginResponse = {
    id: string;
    token: string;
}