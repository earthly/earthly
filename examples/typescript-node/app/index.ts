import express from 'express';
import {sayHello} from './server';

const app = express();

app.get('/hello', (req, res) => {
    const who = req?.query?.who;
    if (who && typeof who !== 'string') {
        res.status(400).send(`provided who is not a string`);
        return;
    }
    const hello = sayHello(who);
    res.status(200).send(hello);
});

app.listen('8080', () => console.log(`Hello world app listening on port 8080`));
