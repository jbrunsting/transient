<template>
  <div class="home">
    <Signup v-if="!authenticated" v-on:signup="updateAuth()"/>
  </div>
</template>

<script>
import Signup from '@/components/Signup.vue';

export default {
    name: 'home',
    props: {
        authenticated: Boolean,
    },
    components: {
        Signup,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
