import {check} from 'k6';
import {request as claimReq} from './requests/claim.js';

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
        userId: `${vId}`.padStart(8, 0),
        amount: data.data.amount,
    };

    const resp = claimReq(input, {baseUrl});

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
