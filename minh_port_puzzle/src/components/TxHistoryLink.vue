<style scoped lang="less">
.link {
  font-size: 0.8em;
  text-align: center;
}
</style>

<template >
  <span>
    <tx-history-panel v-if="isTxPanelOpen" class="tx-history-panel" @close="isTxPanelOpen = false"></tx-history-panel>
    <a class="link" @click="viewTxHistory">View Transactions</a>
  </span>
</template>

<script>
import TxHistoryPanel from "./TxHistoryPanel";

import store from "../store";
export default {
  name: "TutorialPage",
  data() {
    return {
      globalData: store.data,
      isTxPanelOpen: false
    };
  },
  components: {
    TxHistoryPanel
  },
  methods: {
    viewTxHistory() {
      this.isTxPanelOpen = true;
    }
  }
};
</script>
