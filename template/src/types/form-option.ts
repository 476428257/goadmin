export interface FormOption {
    list: FormOptionList[];
    labelWidth?: number | string;
    span?: number;
}

export interface FormOptionList {
    prop: string;
    label: string;
    type: string;
    placeholder?: string;
    disabled?: boolean;
    multiple?: boolean;
    opts?: any[];
    format?: string;
    activeValue?: any;
    inactiveValue?: any;
    suffixLink?: string;
    suffixTooltip?: string;
    activeText?: string;
    inactiveText?: string;
    required?: boolean;
    inputType?: string;
    span?: number; // 新增：单个字段的布局宽度
    height?: string; // 新增：用于设置富文本编辑器高度
    extend?: string; // 新增：扩展配置，JSON字符串格式
}