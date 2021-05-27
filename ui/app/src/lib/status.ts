import {Certificate} from "@/api";
import dayjs from "dayjs";

export default class Status {
    constructor(readonly certificate: Certificate) {
    }

    get soonExpire() {
        const m = this.expiredAt.diff(dayjs(), 'months');
        return m == 0 && !this.expired;
    }

    get expired() {
        return this.expiredAt.isBefore(dayjs())
    }

    get expiredAt() {
        return dayjs(this.certificate.expire_at)
    }

    get duration() {
        return dayjs(this.certificate.expire_at).diff(this.certificate.updated_at!, 'days')
    }
}