import type { Hono } from '../../hono';
import type { Env, Schema } from '../../types';
import type { ALBRequestContext, ApiGatewayRequestContext, ApiGatewayRequestContextV2, Handler, LambdaContext } from './types';
export type LambdaEvent = APIGatewayProxyEvent | APIGatewayProxyEventV2 | ALBProxyEvent;
export interface APIGatewayProxyEventV2 {
    version: string;
    routeKey: string;
    headers: Record<string, string | undefined>;
    multiValueHeaders?: undefined;
    cookies?: string[];
    rawPath: string;
    rawQueryString: string;
    body: string | null;
    isBase64Encoded: boolean;
    requestContext: ApiGatewayRequestContextV2;
    queryStringParameters?: {
        [name: string]: string | undefined;
    };
    pathParameters?: {
        [name: string]: string | undefined;
    };
    stageVariables?: {
        [name: string]: string | undefined;
    };
}
export interface APIGatewayProxyEvent {
    version: string;
    httpMethod: string;
    headers: Record<string, string | undefined>;
    multiValueHeaders?: {
        [headerKey: string]: string[];
    };
    path: string;
    body: string | null;
    isBase64Encoded: boolean;
    queryStringParameters?: Record<string, string | undefined>;
    requestContext: ApiGatewayRequestContext;
    resource: string;
    multiValueQueryStringParameters?: {
        [parameterKey: string]: string[];
    };
    pathParameters?: Record<string, string>;
    stageVariables?: Record<string, string>;
}
export interface ALBProxyEvent {
    httpMethod: string;
    headers?: Record<string, string | undefined>;
    multiValueHeaders?: Record<string, string[] | undefined>;
    path: string;
    body: string | null;
    isBase64Encoded: boolean;
    queryStringParameters?: Record<string, string | undefined>;
    multiValueQueryStringParameters?: {
        [parameterKey: string]: string[];
    };
    requestContext: ALBRequestContext;
}
type WithHeaders = {
    headers: Record<string, string>;
    multiValueHeaders?: undefined;
};
type WithMultiValueHeaders = {
    headers?: undefined;
    multiValueHeaders: Record<string, string[]>;
};
export type APIGatewayProxyResult = {
    statusCode: number;
    statusDescription?: string;
    body: string;
    cookies?: string[];
    isBase64Encoded: boolean;
} & (WithHeaders | WithMultiValueHeaders);
export declare const streamHandle: <E extends Env = Env, S extends Schema = {}, BasePath extends string = "/">(app: Hono<E, S, BasePath>) => Handler;
/**
 * Accepts events from API Gateway/ELB(`APIGatewayProxyEvent`) and directly through Function Url(`APIGatewayProxyEventV2`)
 */
export declare const handle: <E extends Env = Env, S extends Schema = {}, BasePath extends string = "/">(app: Hono<E, S, BasePath>) => (<L extends LambdaEvent>(event: L, lambdaContext?: LambdaContext) => Promise<APIGatewayProxyResult & (L extends {
    multiValueHeaders: Record<string, string[]>;
} ? WithMultiValueHeaders : WithHeaders)>);
export declare abstract class EventProcessor<E extends LambdaEvent> {
    protected abstract getPath(event: E): string;
    protected abstract getMethod(event: E): string;
    protected abstract getQueryString(event: E): string;
    protected abstract getHeaders(event: E): Headers;
    protected abstract getCookies(event: E, headers: Headers): void;
    protected abstract setCookiesToResult(result: APIGatewayProxyResult, cookies: string[]): void;
    createRequest(event: E): Request;
    createResult(event: E, res: Response): Promise<APIGatewayProxyResult>;
    setCookies(event: E, res: Response, result: APIGatewayProxyResult): void;
}
export declare class EventV2Processor extends EventProcessor<APIGatewayProxyEventV2> {
    protected getPath(event: APIGatewayProxyEventV2): string;
    protected getMethod(event: APIGatewayProxyEventV2): string;
    protected getQueryString(event: APIGatewayProxyEventV2): string;
    protected getCookies(event: APIGatewayProxyEventV2, headers: Headers): void;
    protected setCookiesToResult(result: APIGatewayProxyResult, cookies: string[]): void;
    protected getHeaders(event: APIGatewayProxyEventV2): Headers;
}
export declare class EventV1Processor extends EventProcessor<Exclude<LambdaEvent, APIGatewayProxyEventV2>> {
    protected getPath(event: Exclude<LambdaEvent, APIGatewayProxyEventV2>): string;
    protected getMethod(event: Exclude<LambdaEvent, APIGatewayProxyEventV2>): string;
    protected getQueryString(event: Exclude<LambdaEvent, APIGatewayProxyEventV2>): string;
    protected getCookies(event: Exclude<LambdaEvent, APIGatewayProxyEventV2>, headers: Headers): void;
    protected getHeaders(event: APIGatewayProxyEvent): Headers;
    protected setCookiesToResult(result: APIGatewayProxyResult, cookies: string[]): void;
}
export declare class ALBProcessor extends EventProcessor<ALBProxyEvent> {
    protected getHeaders(event: ALBProxyEvent): Headers;
    protected getPath(event: ALBProxyEvent): string;
    protected getMethod(event: ALBProxyEvent): string;
    protected getQueryString(event: ALBProxyEvent): string;
    protected getCookies(event: ALBProxyEvent, headers: Headers): void;
    protected setCookiesToResult(result: APIGatewayProxyResult, cookies: string[]): void;
}
export declare const getProcessor: (event: LambdaEvent) => EventProcessor<LambdaEvent>;
export declare const isContentTypeBinary: (contentType: string) => boolean;
export declare const isContentEncodingBinary: (contentEncoding: string | null) => boolean;
export {};
