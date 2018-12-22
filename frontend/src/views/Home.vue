<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <h2>{{ data }}</h2>
    <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue';

export default {
    name: 'home',
    data() {
        return {
            data: '',
        };
    },
    components: {
        HelloWorld,
    },
    methods: {
        mounted() {
            console.log('HI');
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
