<template>
  <div class="about">
    <h1>Welcome {{ username }}</h1>
    <p>{{ testCall }}</p>
  </div>
</template>
<script>
import Nav from '@/components/Nav.vue';

export default {
    name: 'about',
    data() {
        return {
            username: this.$cookie.get('username'),
            testCall: '',
        };
    },
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
    created() {
        this.$http.get('/api/')
            .then((response) => {
                this.testCall = response.data.message;
            }).catch((e) => {
                console.log(`Error ${JSON.stringify(e)}`);
            });
        this.data = 'test';
    },
};
</script>
