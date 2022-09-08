import axios from 'axios';
import {sayHello} from "../src/server";

describe('sayHello', () => {
    const call = async (who?: string) => {
        const url = who ? `http://localhost:8080/hello?who=${who}` : 'http://localhost:8080/hello';
        const response = await axios.get(url);
        expect(response.status).toBe(200);
        return response.data;
    }

    it('should say Hello Earthly if nothing is passed', async () => {
        expect(await call()).toBe('Hello Earthly');
    });

    it('should say Hello World if World is passed', async () => {
        expect(await call('World')).toBe('Hello World');
    });
});

export {};