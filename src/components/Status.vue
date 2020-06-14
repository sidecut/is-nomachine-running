<template>
  <div>
    <h2 :class="{ error: !connected, header: true }">
      <div>&nbsp;</div>
      <div>
        NoMachine status on {{ hostName }}:
        <span :class="{ hidden: !loading }">
          <img
            src="Spinner-1s-200px.svg"
            alt="loading"
            class="loading-spinner"
          />
        </span>
      </div>
      <div class="setup-gear" @click="settingsClick">⚙️</div>
    </h2>
    <div v-if="initialized">
      <div v-if="connected" style="font-size: 150%;">
        <div v-if="isRunning" class="host-running">
          NoMachine host process is running.
        </div>
        <div v-if="isRunning">
          <div v-if="hasClient" class="attached-client">
            A client IS attached -- use caution when connecting
          </div>
          <div v-else class="no-attached-client">
            No client is attached -- free to connect
          </div>
        </div>
        <div v-else class="host-not-running">
          NoMachine host process is NOT running.
        </div>
      </div>
      <div v-else style="font-size: 150%;">
        <div class="error">Can't get status from {{ hostName }}</div>
      </div>
    </div>
    <button
      v-if="initialized"
      :disabled="loading"
      class="refresh-button"
      @click="refreshClick"
    >
      Refresh
    </button>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";

@Component
export default class Status extends Vue {
  hostName = "";
  initialized = false;
  connected = true;
  isRunning = false;
  hasClient = false;
  loading = false;
  timerHandle = -1;

  created() {
    this.hostName = localStorage.getItem("hostName") ?? "";
  }

  @Watch("hostName")
  storeHostName(newValue: string, oldValue: string) {
    if (newValue) {
      localStorage.setItem("hostName", newValue);
    } else {
      localStorage.removeItem("hostName");
    }

    document.title = `${newValue} NoMachine status`;
  }

  mounted() {
    this.getStatus();
  }

  get apiAddress(): string {
    if (this.hostName) {
      return `//${this.hostName}/api`;
    } else {
      return "./api";
    }
  }

  getStatus() {
    try {
      this.loading = true;
      fetch(this.apiAddress)
        .then((response) => {
          if (response.ok) {
            response.json().then((data: ApiData) => {
              this.connected = true;
              this.hostName = data.HostName;
              console.info(this);
              console.info(
                "data.HostName, this.hostName",
                data.HostName,
                this.hostName
              );
              this.isRunning = data.NoMachineRunning;
              this.hasClient = data.ClientAttached;
            });
          }
        })
        .catch((err) => {
          this.connected = false;
        })
        .finally(() => {
          this.initialized = true;
          this.loading = false;
          this.setUpTimer();
        });
    } catch (err) {
      this.initialized = true;
      this.loading = false;
      this.connected = false;
      this.setUpTimer();
    }
  }
  refreshClick() {
    this.clearTimer();
    this.getStatus();
  }
  clearTimer() {
    window.clearTimeout(this.timerHandle);
    this.timerHandle = -1;
  }
  setUpTimer() {
    this.timerHandle = window.setTimeout(this.getStatus, 45000);
  }

  settingsClick() {
    let newHostName = window.prompt("Enter hostname", this.hostName);
    if (newHostName != null && newHostName != this.hostName) {
      let keepHostname = window.confirm(
        `Use "${newHostName}" instead of "${this.hostName}"?`
      );
      if (keepHostname) {
        this.hostName = newHostName ?? "";
        this.refreshClick();
      }
    }
  }
}

interface ApiData {
  HostName: string;
  NoMachineRunning: boolean;
  ClientAttached: boolean;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
/* body {
  font-family: Arial, Helvetica, sans-serif;
  margin-left: 1rem;
}
div {
  margin-top: 1rem; */
/* font-size: 150%; */
/* } */

.header {
  display: flex;
  justify-content: space-between;
}
.setup-gear {
  /* align-self: flex-end; */
  cursor: pointer;
}

.attached-client {
  font-weight: bold;
  color: coral;
}
.no-attached-client {
  color: green;
}
.host-not-running {
  font-weight: bold;
  color: red;
}
.host-running {
  font-weight: bold;
  color: green;
}
.error {
  color: red;
}
.loading-spinner {
  height: 28pt;
  vertical-align: bottom;
}
.refresh-button {
  margin-top: 2rem;
  height: 2rem;
  width: 4rem;
}
.hidden {
  visibility: hidden;
}
</style>
