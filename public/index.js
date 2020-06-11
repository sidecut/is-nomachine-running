var app = new Vue({
  el: "#app",
  data: {
    connected: true,
    isRunning: false,
    hasClient: false,
    loading: false,
  },
  mounted: function () {
    window.setInterval(this.getStatus, 6000);
    this.isRunning = !!document.getElementById("is-running").nodeValue;
    this.hasClient = !!document.getElementById("client-attached").nodeValue;
  },
  methods: {
    getStatus: function () {
      try {
        this.loading = true;
        fetch("./api")
          .then((response) => {
            this.loading = false;
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
            this.loading = false;
          });
      } catch (err) {
        this.loading = false;
        this.connected = false;
      }
    },
  },
});
