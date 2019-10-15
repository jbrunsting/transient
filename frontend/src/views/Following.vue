<template>
  <div class="following">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <FollowingRecommends v-on:follow="updateFollowings()"/>
      <Followings v-bind:followings="followings"  v-on:follow="updateFollowings()"/>
      <Follow v-bind:followings="followings" v-on:follow="updateFollowings()"/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Followings from '@/components/Followings.vue';
import FollowingRecommends from '@/components/FollowingRecommends.vue';
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
        FollowingRecommends,
        Follow,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        updateFollowings() {
            this.$http.getProtected('/api/followings')
                .then((response) => {
                    this.followings = response.data;
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
