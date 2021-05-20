import {check} from 'k6';
import {request as usageReq} from './requests/usage.js';

const config = JSON.parse(open('./config.json'));
const data = JSON.parse(open('./data.json'));

const vus = data.vus
const iterations = data.iterations * vus;

export let options = {
    vus,
    iterations,
};

export default function () {
    const vId = __VU;

    const baseUrl = config.hosts[vId % config.hosts.length];
    const input = {
        id: data.data.id,
        amount: data.data.amount,
    };

    const resp = usageReq(input, {baseUrl});

    check(resp, {
        'success': (r) => {
            return r.status === 200;
        },
    });

    check(resp, {
        'error': (r) => {
            return r.status !== 200;
        },
    });
}
