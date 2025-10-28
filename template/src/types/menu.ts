export interface Menus {
    id: string;
    pid?: string;
    icon?: string;
    route: string;
    title: string;
    is_menu: number;
    status: number;
    children?: Menus[];
}