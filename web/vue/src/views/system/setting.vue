<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <strong> 系统设置 </strong>
      </div>
    </template>
    <el-form ref="form" :model="setting" label-width="100px" style="width: 500px;">
      <el-form-item label="Logo" prop="logo">
        <el-input v-model="setting.logo"></el-input>
      </el-form-item>
      <el-form-item label="网站名称" prop="title">
        <el-input v-model="setting.title"></el-input>
      </el-form-item>
      <el-form-item style="text-align: right;">
        <el-button type="primary" @click="submit()">保存</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import systemService from '@/api/system'

export default {
  name: "system-setting",
  data() {
    return {
      setting: {}
    }
  },
  created() {
    let _this = this
    systemService.getSystemSetting(function (data) {
      _this.setting = data
    })
  },
  methods: {
    submit() {
      let _this = this
      systemService.updateSystemSetting(_this.setting, function () {
        _this.$message.success("操作成功")
        _this.$store.commit('setSystemSetting', _this.setting)
      })
    }
  }
}
</script>

<style scoped>

</style>