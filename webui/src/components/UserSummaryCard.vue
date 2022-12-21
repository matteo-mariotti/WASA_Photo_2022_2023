<script>

import backend from "@/services/axios";

export default {
  props: ["username"],
  data: function () {
    return {
    }
  },
  methods: {
  show(){
    this.$router.push(`/users/${this.username}`)
  }
  },
}

</script>


<template>

  <div class="card m-2">
    <div class="card-body">
      <div class="d-flex">
        <h6>{{this.username}}</h6>
        <button @click="show()" class="btn btn-outline-success btn-sm ms-auto"> View profile
        </button>
      </div>
    </div>
  </div>

</template>

<style>
</style>