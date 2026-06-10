import "vue-class-component/hooks"; // import hooks type to enable auto-complete
import Vue from 'vue';
import App from './App.vue';
import './registerServiceWorker';

Vue.config.productionTip = false;

new Vue({
  render: (h) => h(App),
}).$mount('#app');
