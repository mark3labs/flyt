type PreviewParameters = {
    outDir?: string;
};
export declare function preview({ outDir }?: PreviewParameters): Promise<import("@hono/node-server").ServerType & {
    port: number;
}>;
export {};
//# sourceMappingURL=preview.d.js.map