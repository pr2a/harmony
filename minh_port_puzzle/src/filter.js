import Vue from 'vue'
import moment from "moment";
export function formatTimestamp(timestamp) {
    return moment(timestamp).format('MM/DD/YYYY hh:mm:ss');
}

Vue.filter('timestamp', formatTimestamp);

export function shortenHash(hash) {
    if (!hash || hash.length <= 10) return hash;
    return hash.substr(0, 5) + "..." + hash.substr(hash.length - 5);
}
Vue.filter('shorten', shortenHash);
