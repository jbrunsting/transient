<template>
  <div class="following">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <Followings v-bind:followings="followings"/>
      <Follow v-bind:followings="followings" v-on:follow="updateFollowings()"/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Followings from '@/components/Followings.vue';
import Follow from '@/components/Follow.vue';

export default {
    name: 'following',
    data() {
        return {
            followings: [],
        };
    },
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Followings,
        Follow,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        updateFollowings() {
            this.$http.getProtected('/api/followings')
                .then((response) => {
                    console.log("GOT response")
                    this.followings = response.data;
                    console.log("Folloowings are " + this.followings)
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
    created() {
        this.updateFollowings();
    },
};
</script>
