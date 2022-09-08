import {sayHello} from './server';

describe('sayHello', () => {
    it('should say Hello Earthly if nothing is passed', () => {
        expect(sayHello()).toBe('Hello Earthly');
    });

    it('should say Hello World if World is passed', () => {
        expect(sayHello('World')).toBe('Hello World');
    });
})