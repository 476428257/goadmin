
export interface Role {
    id: number;
    title: string;
    status: boolean;
    permiss: string[]
    rule: any
    children?: Role[];
}