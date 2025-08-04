/**
 * @module
 * JSON Web Token (JWT)
 * https://datatracker.ietf.org/doc/html/rfc7519
 */
import type { SignatureAlgorithm } from './jwa';
import type { HonoJsonWebKey, SignatureKey } from './jws';
import type { JWTPayload } from './types';
export interface TokenHeader {
    alg: SignatureAlgorithm;
    typ?: "JWT";
    kid?: string;
}
export declare function isTokenHeader(obj: unknown): obj is TokenHeader;
export declare const sign: (payload: JWTPayload, privateKey: SignatureKey, alg?: SignatureAlgorithm) => Promise<string>;
export declare const verify: (token: string, publicKey: SignatureKey, alg?: SignatureAlgorithm) => Promise<JWTPayload>;
export declare const verifyFromJwks: (token: string, options: {
    keys?: HonoJsonWebKey[];
    jwks_uri?: string;
}, init?: RequestInit) => Promise<JWTPayload>;
export declare const decode: (token: string) => {
    header: TokenHeader;
    payload: JWTPayload;
};
export declare const decodeHeader: (token: string) => TokenHeader;
