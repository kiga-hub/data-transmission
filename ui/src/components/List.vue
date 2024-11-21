<template>
<div class="layout">
  <Layout>
    <Header style="height:56px">
    </Header>
    <Content style="padding: 24px 50px">
      <Card>
        <p slot="title">
          <Icon type="ios-list"></Icon>
          部署任务列表
        </p>
        <div slot="extra">
          <Button icon="md-add" type="primary" @click="onCreateTask">创建</Button>
          <Button icon="md-arrow-forward" type="primary" @click="onUpgrade">升级数据模型</Button>
        </div>
        <Table :columns="columns" :data="datas" max-height="1000" border></Table>
      </Card>
    </Content>
    <Footer class="layout-footer-center"></Footer>
  </Layout>
  <Modal
    v-model="createModal"
    title="创建部署任务"
    width="1024"
    @on-ok="onClickRun"
    >

    <Row type="flex" justify="center" :gutter="32">
      <Col style="height: 512px; margin-top: 64px;">
        <Form :model="param.extra_vars.basic" label-position="left" :label-width="100" ref="basic" :rules="basicRuleValidate">
          <FormItem label="项目名称" prop="project">
            <Input v-model="param.extra_vars.basic.project"></Input>
          </FormItem>
          <FormItem label="用途">
            <Select v-model="param.extra_vars.basic.use" style="width:260px">
              <Option value="dev">开发</Option>
              <Option value="test">测试</Option>
              <Option value="release">发布</Option>
            </Select>
          </FormItem>
          <FormItem label="部署类型">
            <Select v-model="param.extra_vars.basic.tag" style="width:260px">
              <Option value="docker">Docker</Option>
              <Option value="rancher">Rancher</Option>
              <Option value="iot">物联网平台</Option>
              <Option value="voiceprint">工业声纹状态监测系统</Option>
              <Option value="app">应用大屏</Option>
            </Select>
          </FormItem>
          <FormItem label="应用" v-if="param.extra_vars.basic.tag=='app'">
            <Select v-model="param.extra_vars.basic.app" style="width:260px">
              <Option value="xldsj">溪洛渡水机</Option>
            </Select>
          </FormItem>
        </Form>
        <Form :model="param.extra_vars.machine" label-position="left" :label-width="100" ref="machine" :rules="machineRuleValidate">
          <FormItem label="服务器地址" prop="ip">
            <Input v-model="param.extra_vars.machine.ip" @on-change="onChangeip"></Input>
          </FormItem>
          <FormItem label="主机名" prop="hostname">
            <Input v-model="param.extra_vars.machine.hostname"></Input>
          </FormItem>
          <FormItem label="临时文件目录">
            <Input v-model="param.extra_vars.machine.dest_dir" disabled></Input>
          </FormItem>
          <FormItem label="用户名">
            <Input v-model="param.extra_vars.machine.user" disabled></Input>
          </FormItem>
          <FormItem label="密码" prop="password">
            <Input v-model="param.extra_vars.machine.password"></Input>
          </FormItem>
        </Form>
      </Col>
      <Col style="height: 512px; margin-top: 64px;">
        <Form :model="param.extra_vars.install" label-position="left" :label-width="140">
          <FormItem label="系统服务安装">
            <CheckboxGroup v-model="param.extra_vars.install">
              <Checkbox label="NTP"></Checkbox>
            </CheckboxGroup>
          </FormItem>
        </Form>
        <Form :model="param.extra_vars.init" label-position="left" :label-width="140">
          <FormItem label="中间件初始化">
            <CheckboxGroup v-model="param.extra_vars.init">
              <Checkbox label="Mysql"></Checkbox>
              <Checkbox label="Pulsar"></Checkbox>
              <Checkbox label="Clickhouse"></Checkbox>
              <Checkbox label="Iotdb"></Checkbox>
              <Checkbox label="Redis"></Checkbox>
            </CheckboxGroup>
          </FormItem>
        </Form>
        <Form :model="param.extra_vars.rancher" label-position="left" :label-width="140" ref="rancher" :rules="rancherRuleValidate">
          <FormItem label="Rancher连接地址">
            <Input v-model="param.extra_vars.rancher.url" disabled></Input>
          </FormItem>
          <FormItem label="Rancher登录密码" prop="password">
            <Input v-model="param.extra_vars.rancher.password"></Input>
          </FormItem>
          <FormItem label="集群名称" prop="cluster_name">
            <Input v-model="param.extra_vars.rancher.cluster_name"></Input>
          </FormItem>
          <FormItem label="项目名称" prop="project_name">
            <Input v-model="param.extra_vars.rancher.project_name"></Input>
          </FormItem>
          <FormItem label="命名空间" prop="namespace">
            <Input v-model="param.extra_vars.rancher.namespace"></Input>
          </FormItem>
        </Form>
      </Col>
    </Row>
  </Modal>
  <Modal
    v-model="detailModal"
    :title="detailTitle"
    width="1024"
    footer-hide>
    <Row type="flex" justify="center" :gutter="32" style="margin: 20px;">
      <Col>
        开始时间:
        <Tag color="blue"> {{ detail.StartTime }} </Tag>
      </Col>
      <Col>
        结束时间:
        <Tag color="blue"> {{ detail.EndTime }} </Tag>
      </Col>
      <Col>
        状态:
        <Tag :color="getTaskColor(detail.State)">{{ getTaskState(detail.State) }}</Tag>
      </Col>
    </Row>
    <Tabs value="log">
      <TabPane label="日志" name="log">
        <Card v-if="detail.Log">
          <div class="markdown-container" id="log_div">
            <markdown-it-vue class="markdown-body" :content="detail.Log" :options="options" />
          </div>
        </Card>
      </TabPane>
      <TabPane label="配置" name="config">
        <Card v-if="detail.ExtraVars">
          <div class="markdown-container">
            <json-viewer :value="detail.ExtraVars" class="markdown-body"></json-viewer>
          </div>
        </Card>
      </TabPane>
    </Tabs>
  </Modal>
</div>
</template>

<script>
import MarkdownItVue from 'markdown-it-vue'
import 'markdown-it-vue/dist/markdown-it-vue.css'
import {
  getList, deleteTask, getTask, runAnsible
} from '@/api/data'

export default {
  name: 'Create',
  components: {
    MarkdownItVue
  },
  data () {
    return {
      createModal: false,
      param: {
        playbook: './playbook/main.yml',
        extra_vars: {
          basic: {
            project: '',
            tag: 'docker', // app
            use: 'dev', // test、release
            app: 'xldsj'
          },
          machine: {
            ip: '',
            password: '',
            user: 'root',
            hostname: '',
            dest_dir: '/root'
          },
          install: ['NTP'],
          init: ['Mysql', 'Pulsar', 'Clickhouse', 'Iotdb', 'Redis'],
          rancher: {
            url: 'https://127.0.0.1:8443',
            password: '123456',
            cluster_name: 'aithu',
            project_name: 'project',
            namespace: 'middle'
          }
        }
      },
      basicRuleValidate: {
        project: [
          { required: true, message: '项目名称不能为空', trigger: 'blur' }
        ]
      },
      machineRuleValidate: {
        ip: [
          { required: true, message: '服务器地址不能为空', trigger: 'blur' }
        ],
        hostname: [
          { required: true, message: '服务器名称不能为空', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '密码不能为空', trigger: 'blur' }
        ]
      },
      rancherRuleValidate: {
        password: [
          { required: true, message: 'Rancher密码不能为空', trigger: 'blur' }
        ],
        cluster_name: [
          { required: true, message: '集群名称不能为空', trigger: 'blur' }
        ],
        project_name: [
          { required: true, message: '项目名称不能为空', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '命名空间不能为空', trigger: 'blur' }
        ]
      },
      columns: [{
        title: '项目',
        key: 'Project'
      }, {
        title: '用途',
        key: 'Use',
        render: (h, params) => {
          let text = '开发'
          if (params.row.Use === 'test') {
            text = '测试'
          } else if (params.row.Use === 'release') {
            text = '发布'
          }
          return h('div', [
            h('span', {
            }, text)
          ])
        }
      }, {
        title: '类型标签',
        key: 'Tag',
        render: (h, params) => {
          let text = ''
          if (params.row.Tag === 'docker') {
            text = 'Docker'
          } else if (params.row.Tag === 'rancher') {
            text = 'Rancher'
          } else if (params.row.Tag === 'iot') {
            text = '物联网平台'
          } else if (params.row.Tag === 'voiceprint') {
            text = '工业声纹状态监测系统'
          } else if (params.row.Tag === 'app') {
            text = '应用大屏-' + params.row.App
          }
          return h('div', [
            h('span', {
            }, text)
          ])
        }
      }, {
        title: 'IP',
        key: 'IP'
      }, {
        title: '服务器名',
        key: 'HostName'
      }, {
        title: '状态',
        key: 'State',
        render: (h, params) => {
          var text = '完成'
          var color = 'blue'
          if (params.row.State === 1) {
            text = '运行中'
            color = 'green'
          } else if (params.row.State === 3) {
            text = '失败'
            color = 'red'
          }
          return h('div', [
            h('Tag', {
              props: {
                color: color
              }
            }, text)
          ])
        }
      }, {
        title: '开始时间',
        key: 'StartTime'
      }, {
        title: '结束时间',
        key: 'EndTime'
      }, {
        title: '操作',
        key: 'action',
        width: 150,
        align: 'center',
        render: (h, params) => {
          return h('div', [
            h('Button', {
              props: {
                ghost: true,
                size: 'small',
                type: 'info'
              },
              style: {
                marginRight: '5px'
              },
              on: {
                click: () => {
                  this.detailModal = true
                  this.getTask(params.row.Date)
                }
              }
            }, '查看'),
            h('Button', {
              props: {
                type: 'error',
                ghost: true,
                disabled: params.row.State === 1,
                size: 'small'
              },
              on: {
                click: () => {
                  this.$Modal.confirm({
                    title: '删除特征',
                    onOk: () => {
                      this.deleteTask(params.row.Date)
                    }
                  })
                }
              }
            }, '删除')
          ])
        }
      }],
      datas: [],
      detailModal: false,
      detailTitle: '',
      detail: {},
      options: {
        markdownIt: {
          linkify: true,
          breaks: true
        },
        linkAttributes: {
          attrs: {
            target: '_blank',
            rel: 'noopener'
          }
        }
      },
      host: '',
      isDevelopment: process.env.NODE_ENV !== 'production'
    }
  },
  created () {
    if (this.isDevelopment) {
      this.host = '192.168.8.120'
    } else {
      this.host = location.host
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    onUpgrade () {
      this.$router.push({ path: '/upgrade' })
    },
    onChangeip (e) {
      this.param.extra_vars.rancher.url = 'https://' + this.param.extra_vars.machine.ip + ':8443'
    },
    createTask () {
      var _this = this
      runAnsible({
        playbook: this.param.playbook,
        extra_vars: {
          basic: this.param.extra_vars.basic,
          machine: this.param.extra_vars.machine,
          install: {
            ntp: this.param.extra_vars.install.indexOf('NTP') !== -1
          },
          init: {
            mysql: this.param.extra_vars.init.indexOf('Mysql') !== -1,
            pulsar: this.param.extra_vars.init.indexOf('Pulsar') !== -1,
            clickhouse: this.param.extra_vars.init.indexOf('Clickhouse') !== -1,
            iotdb: this.param.extra_vars.init.indexOf('Iotdb') !== -1,
            redis: this.param.extra_vars.init.indexOf('Redis') !== -1
          },
          rancher: this.param.extra_vars.rancher
        }
      }).then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        _this.$Message.info(res.data.msg)
        _this.getList()
      })
    },
    onClickRun () {
      this.$refs['basic'].validate((valid) => {
        if (!valid) {
          this.$Message.error('基础信息填写错误')
          return
        }
        this.$refs['machine'].validate((valid) => {
          if (!valid) {
            this.$Message.error('服务器信息填写错误')
            return
          }
          this.$refs['rancher'].validate((valid) => {
            if (!valid) {
              this.$Message.error('Rancher信息填写错误')
            }
            this.createTask()
          })
        })
      })
    },
    onCreateTask () {
      this.createModal = true
    },
    getTaskState (state) {
      if (state === 1) {
        return '运行中'
      } else if (state === 2) {
        return '完成'
      } else if (state === 3) {
        return '失败'
      }
    },
    getTaskColor (state) {
      if (state === 1) {
        return 'info'
      } else if (state === 2) {
        return 'success'
      } else if (state === 3) {
        return 'error'
      }
    },
    getTask (date) {
      if (!this.detailModal) {
        return
      }
      var _this = this
      getTask({date: date}).then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        _this.detail = res.data.data
        _this.detailTitle = '任务 - ' + date
        if (_this.detail.State === 1) {
          setTimeout(function () {
            _this.getTask(date)
          }, 2000)
        }
      })
    },
    deleteTask (date) {
      var _this = this
      deleteTask({date: date}).then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        _this.datas = _this.datas.filter(item => {
          return item.Date !== date
        })
        _this.$Message.info(res.data.msg)
      })
    },
    getList () {
      var _this = this
      getList().then(res => {
        if (res.data.code !== 200) {
          _this.$Message.error(res.data.msg)
          return
        }
        _this.datas = res.data.data
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.layout{
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
}
.layout-footer-center{
    text-align: center;
}
.markdown-container {
  max-height: 500px; /* 设置最大高度，确保内容超出时显示滚动条 */
  overflow-y: auto; /* 设置为自动显示滚动条 */
}

.markdown-body {
  padding: 15px; /* 可选：为Markdown内容添加内边距 */
}
</style>
