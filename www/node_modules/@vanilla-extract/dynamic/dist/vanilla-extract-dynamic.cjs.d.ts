type CSSVarFunction = `var(--${string})`;
type Contract = {
    [key: string]: CSSVarFunction | null | Contract;
};
type Primitive = string | boolean | number | null | undefined;
type MapLeafNodes<Obj, LeafType> = {
    [Prop in keyof Obj]: Obj[Prop] extends Primitive ? LeafType : Obj[Prop] extends Record<string | number, any> ? MapLeafNodes<Obj[Prop], LeafType> : never;
};

type Styles = {
    [cssVarName: string]: string;
};
declare function assignInlineVars(vars: Record<string, string | undefined | null>): Styles;
declare function assignInlineVars<ThemeContract extends Contract>(contract: ThemeContract, tokens: MapLeafNodes<ThemeContract, string>): Styles;

declare function setElementVars(element: HTMLElement, vars: Record<string, string | undefined | null>): void;
declare function setElementVars<ThemeContract extends Contract>(element: HTMLElement, contract: ThemeContract, tokens: MapLeafNodes<ThemeContract, string>): void;

export { assignInlineVars, setElementVars };
