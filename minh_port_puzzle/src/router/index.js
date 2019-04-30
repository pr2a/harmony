import Puzzle from '@/components/Puzzle';
import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Puzzle',
      component: Puzzle
    }
  ]
});
