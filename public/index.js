var app = new Vue({
  el: "#app",
  data: {
    initialized: false,
    connected: true,
    isRunning: false,
    hasClient: false,
    loading: false,
    timerHandle: -1,
  },
  mounted: function () {
    this.getStatus();
  },
  methods: {
    getStatus: function () {
      try {
        this.loading = true;
        fetch("./api")
          .then((response) => {
            if (response.ok) {
              response.json().then((data) => {
                this.connected = true;
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
    },
    refreshClick() {
      this.clearTimer();
      this.getStatus();
    },
    clearTimer() {
      window.clearTimeout(this.timerHandle);
      this.timerHandle = -1;
    },
    setUpTimer() {
      this.timerHandle = window.setTimeout(this.getStatus, 45000);
    },
  },
});
