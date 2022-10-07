export function sayHello(who?: string): string {
    who = who ?? 'Earthly';
    return `Hello ${who}`;
}
