
export interface User {
    id: number;
    username: string;
    nickname: string;
    password: string;
    email: string;
    phone: string;
    created_at:string
    updated_at:string
    role: string;
    status:number;
}

export interface Register {
    username: string;
    password: string;
    email: string;
}