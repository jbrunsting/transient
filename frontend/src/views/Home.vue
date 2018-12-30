<template>
  <div class="home">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <Signup v-if="!authenticated" v-on:signup="updateAuth()"/>
    <img alt="Vue logo" src="../assets/logo.png">
    <h2>{{ data }}</h2>
    <h2>{{ authenticated }}</h2>
    <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
import HelloWorld from '@/components/HelloWorld.vue';
import Nav from '@/components/Nav.vue';
import Signup from '@/components/Signup.vue';

export default {
    name: 'home',
    props: {
        authenticated: Boolean,
    },
    data() {
        return {
            data: '',
        };
    },
    components: {
        HelloWorld,
        Nav,
        Signup,
    },
    methods: {
        mounted() {
            console.log('HI');
        },
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.get('/api/')
            .then((response) => {
                this.data = response.data.message;
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
        this.data = 'test';
    },
};
</script>
