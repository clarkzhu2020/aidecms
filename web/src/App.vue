<template>
  <div id="app">
    <NavBar />
    <RouterView />
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'
import Footer from './components/Footer.vue'
import { onMounted } from 'vue'
import AOS from 'aos'
import 'aos/dist/aos.css'

onMounted(() => {
  AOS.init({
    duration: 800,
    once: true,
    offset: 100,
  })
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background: #ffffff;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

:root {
  --primary-color: #667eea;
  --secondary-color: #764ba2;
  --accent-color: #00d4ff;
  --dark-bg: #0f172a;
  --light-bg: #f8fafc;
  --text-dark: #1e293b;
  --text-light: #64748b;
  --gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --gradient-secondary: linear-gradient(135deg, #00d4ff 0%, #0099ff 100%);
  --gradient-dark: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

@media (max-width: 768px) {
  .container {
    padding: 0 15px;
  }
}
</style>
