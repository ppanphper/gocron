<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <strong>{{ form.id > 0 ? '编辑' : '新增' }}</strong>
      </div>
    </template>
    <el-form ref="form" :model="form" :rules="formRules" label-width="180px">
      <el-input v-model="form.id" type="hidden"></el-input>
      <el-row>
        <el-col :span="12">
          <el-form-item label="进程名称" prop="name">
            <el-input v-model.trim="form.name"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="标签">
            <el-input v-model.trim="form.tag" placeholder="通过标签将进程分组"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24" style="margin-bottom: 15px">
          <el-alert title="进程数量最小值为1;" type="info" :closable="false"></el-alert>
          <el-alert title="日志文件为空时,不保存进程文件,如果指定了日志文件,系统会先创建对应的文件，再将进程的输出写入日志文件;" type="info" :closable="false"></el-alert>
          <el-alert title="如果文件已存在,则直接写入，请确保运行节点有保存日志文件的目录的读写权限" type="info" :closable="false"></el-alert>
        </el-col>

        <el-col :span="8">
          <el-form-item label="进程数量" prop="num_proc">
            <el-input v-model="form.num_proc"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="16">
          <el-form-item label="日志文件">
            <el-input v-model.trim="form.log_file" placeholder="日志文件,为空表示不保存日志文件"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <el-form-item label="项目">
            <el-select v-model="form.project_id" prop="project_id" clearable>
              <el-option v-for="project in projects" :value="project.id" :key="project.id" :label="project.name"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="24" style="margin-bottom: 15px">
          <el-alert title="如果指定了运行节点,进程会在指定节点上运行;    如果没有指定节点,进程会在项目关联的节点上运行" type="info" :closable="false"></el-alert>
        </el-col>
        <el-col :span="24">
          <el-form-item label="运行节点">
            <el-select v-model="processHostIds" filterable multiple placeholder="请选择" style="width: 95%" clearable>
              <el-option
                  v-for="item in hosts"
                  :key="item.id"
                  :label="item.alias + ' - ' + item.name"
                  :value="item.id">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="16">
          <el-form-item label="命令" prop="command">
            <el-input
                type="textarea"
                :rows="5"
                size="medium"
                width="100"
                v-model="form.command">
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item>
        <el-button type="primary" @click="submit">保存</el-button>
        <el-button @click="cancel">取消</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import processService from '../../api/process'
import hostService from '../../api/host'
import projectService from "@/api/project";

export default {
  name: 'process-edit',
  data() {
    return {
      form: {
        num_proc: 1,
        log_file: ''
      },
      processHostIds: [],
      formRules: {},
      hosts: [],
      projects:[],
    }
  },
  created() {
    const id = this.$route.params.id
    let _this = this
    if (id) {
      processService.get(id, function (data) {
        _this.form = Object.assign({}, data)
        !!data.hosts && data.hosts.forEach(host => {
          _this.processHostIds.push(host.id)
        })
      })
    }
    hostService.all({}, (hosts) => {
      _this.hosts = hosts
    })
    projectService.all((data) => {
      this.projects = data.projects
    })
  },
  methods: {
    submit() {
      let _this = this
      _this.form.host_ids = _this.processHostIds.join(',')
      processService.store(_this.form, function (resp) {
        console.log(resp)
        _this.$message.success('Success')
        _this.$router.push('/process/index')
      })
    },
    cancel() {
      this.$router.push('/process/index')
    }
  }
}
</script>

<style scoped>

</style>
