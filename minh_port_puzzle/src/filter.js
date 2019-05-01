import Vue from 'vue'
import moment from "moment";
export function formatTimestamp(timestamp) {
    return moment(timestamp).format('MM/DD/YYYY hh:mm:ss');
}

Vue.filter('timestamp', formatTimestamp);
