/**
 * @module
 * JWT utility.
 */
export declare const Jwt: {
    sign: (payload: import("./types").JWTPayload, privateKey: import("./jws").SignatureKey, alg?: import("./jwa").SignatureAlgorithm) => Promise<string>;
    verify: (token: string, publicKey: import("./jws").SignatureKey, alg?: import("./jwa").SignatureAlgorithm) => Promise<import("./types").JWTPayload>;
    decode: (token: string) => {
        header: import("./jwt").TokenHeader;
        payload: import("./types").JWTPayload;
    };
    verifyFromJwks: (token: string, options: {
        keys?: import("./jws").HonoJsonWebKey[];
        jwks_uri?: string;
    }, init?: RequestInit) => Promise<import("./types").JWTPayload>;
};
