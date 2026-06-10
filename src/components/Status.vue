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
      <div class="setup-gear" @click="settingsClick" tabindex="0">⚙️</div>
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

<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";

const hostName = ref(localStorage.getItem("hostName") ?? "");
const initialized = ref(false);
const connected = ref(true);
const isRunning = ref(false);
const hasClient = ref(false);
const loading = ref(false);
const timerHandle = ref(-1);

watch(hostName, (newValue: string) => {
  if (newValue) {
    localStorage.setItem("hostName", newValue);
  } else {
    localStorage.removeItem("hostName");
  }
  document.title = `${newValue} NoMachine status`;
});

const apiAddress = computed((): string => {
  if (hostName.value) {
    if (hostName.value.includes("//")) {
      return `${hostName.value}/api`;
    } else {
      return `//${hostName.value}/api`;
    }
  } else {
    return "./api";
  }
});

function getStatus() {
  try {
    loading.value = true;
    fetch(apiAddress.value)
      .then((response) => {
        if (response.ok) {
          response.json().then((data: ApiData) => {
            connected.value = true;
            if (!hostName.value) {
              hostName.value = data.host_name;
            }
            isRunning.value = data.no_machine_running;
            hasClient.value = data.client_attached;
          });
        }
      })
      .catch(() => {
        connected.value = false;
      })
      .finally(() => {
        initialized.value = true;
        loading.value = false;
        setUpTimer();
      });
  } catch (err) {
    initialized.value = true;
    loading.value = false;
    connected.value = false;
    setUpTimer();
  }
}

function refreshClick() {
  clearTimer();
  getStatus();
}

function clearTimer() {
  window.clearTimeout(timerHandle.value);
  timerHandle.value = -1;
}

function setUpTimer() {
  timerHandle.value = window.setTimeout(getStatus, 15000);
}

function settingsClick() {
  const newHostName = window.prompt(
    "Enter hostname, with optional leading http:// or https:// and optional port",
    hostName.value
  );
  if (newHostName != null && newHostName != hostName.value) {
    const useNewHostname = window.confirm(
      `Use "${newHostName}" instead of "${hostName.value}"?`
    );
    if (useNewHostname) {
      hostName.value = newHostName ?? "";
      refreshClick();
    }
  }
}

onMounted(() => {
  getStatus();
});

interface ApiData {
  host_name: string;
  no_machine_running: boolean;
  client_attached: boolean;
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
