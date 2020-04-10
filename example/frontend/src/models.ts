/* Do not change, this code is generated from Golang structs */


export interface APIError {
    reason: string;
}
export function createAPIErrorFrom(source: any): APIError {
    if ('string' === typeof source) source = JSON.parse(source);
    const result: any = {};
    result.reason = source["reason"];
    return result as APIError;
}
export interface Policy {
    resource: string;
    action: string;
}
export function createPolicyFrom(source: any): Policy {
    if ('string' === typeof source) source = JSON.parse(source);
    const result: any = {};
    result.resource = source["resource"];
    result.action = source["action"];
    return result as Policy;
}

export interface User {
    id: string;
    email: string;
    createdAt: Date;
    updatedAt: Date;
}
export function createUserFrom(source: any): User {
    if ('string' === typeof source) source = JSON.parse(source);
    const result: any = {};
    result.id = source["id"];
    result.email = source["email"];
    result.createdAt = new Date(source["createdAt"]);
    result.updatedAt = new Date(source["updatedAt"]);
    return result as User;
}

export interface ProfileResponse {
    user: User;
    policies: Policy[];
}
export function createProfileResponseFrom(source: any): ProfileResponse {
    if ('string' === typeof source) source = JSON.parse(source);
    const result: any = {};
    result.user = createUserFrom(source["user"]);
    result.policies = source["policies"].map(function(element: any) { return createPolicyFrom(element); });
    return result as ProfileResponse;
}