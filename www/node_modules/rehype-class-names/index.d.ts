import { Plugin } from 'unified';
import { Root } from 'hast';
import { Conditional } from 'hast-util-classnames';
export interface Options {
    [selector: string]: Conditional;
}
declare const rehypeClassNames: Plugin<[(Options | null | undefined)?], Root>;
export default rehypeClassNames;
