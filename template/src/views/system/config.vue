<template>
    <div class="container">
        <div class="panel panel-default panel-intro">
            <div class="panel-heading">
                <ul class="nav nav-tabs">
                    <li v-for="g in groups" :key="g" :class="{ active: currentGroup === g }">
                        <a href="javascript:;" @click="currentGroup = g; onGroupChange()">{{ g }}</a>
                    </li>
                </ul>
            </div>
            <div class="panel-body">
        <el-form ref="formRef" :model="form" label-width="120px" :label-position="labelPosition">
            <el-row :gutter="50">
                <el-col :span="12">
                    <template v-if="items.length">
                        <el-form-item v-for="item in items" :key="item.name" :label="item.title">
                            <template v-if="item.type === 'string'">
                                <el-input v-model="form[item.name]" />
                            </template>
                            <template v-else-if="item.type === 'number'">
                                <el-input-number v-model="form[item.name]" :min="Number.MIN_SAFE_INTEGER" :max="Number.MAX_SAFE_INTEGER" />
                            </template>
                            <template v-else-if="item.type === 'select'">
                                <el-select v-model="form[item.name]" placeholder="请选择">
                                    <el-option v-for="opt in parseSelectOptions(item.content)" :key="opt.value" :label="opt.label" :value="opt.value" />
                                </el-select>
                            </template>
                            <template v-else-if="item.type === 'date'">
                                <el-date-picker 
                                    :type="parseExtend(item.extend).type || 'date'"
                                    v-model="form[item.name]"
                                    :value-format="parseExtend(item.extend).format || 'YYYY-MM-DD'"
                                    placeholder="请选择日期"
                                />
                            </template>
                            <template  v-else-if="item.type === 'daterange'" >
                                <el-date-picker style="width: 400px;"
                                    :type="parseExtend(item.extend).type || 'daterange'"
                                    v-model="form[item.name]"
                                    :value-format="parseExtend(item.extend).format || 'YYYY-MM-DD'"
                                    range-separator="至" 
                                    start-placeholder="开始日期" 
                                    end-placeholder="结束日期"
                                    clearable
                                />
                            </template>
                            <template v-else-if="item.type === 'switch'">
                                <el-switch 
                                    v-model="form[item.name]"
                                    :active-value="parseExtend(item.extend).activeValue || 1" 
                                    :inactive-value="parseExtend(item.extend).inactiveValue || 0"
                                    :active-text="parseExtend(item.extend).activeText || '开启'" 
                                    :inactive-text="parseExtend(item.extend).inactiveText || '关闭'"
                                />
                            </template>
                            <template v-else-if="item.type === 'upload'">
                                <el-upload 
                                    class="upload-container" 
                                    style="border: 1px dashed #ccc; padding: 10px; border-radius: 6px;"
                                    :multiple="parseExtend(item.extend).multiple || false"
                                    :auto-upload="false"
                                    :show-file-list="false" 
                                    :disabled="parseExtend(item.extend).multiple && getDisplayFileList(item.name, item).length > 0"
                                    :on-change="(file, files) => handleFilesChange(files, item.name, item)"
                                    :ref="el => setMainUploadRef(item.name, el)"
                                >
                                    <!-- 多文件预览 -->
                                    <div v-if="parseExtend(item.extend).multiple && getDisplayFileList(item.name, item).length" class="multiple-files">
                                        <div v-for="(f, index) in getDisplayFileList(item.name, item)" :key="index" class="file-item">
                                            <div class="file-preview-container">
                                                <img v-if="isImageFile(f)" :src="f.url" class="file-preview" />
                                                <img v-else src="/image/doc.png" class="file-preview" />
                                                <el-button 
                                                    type="danger" 
                                                    size="small" 
                                                    circle 
                                                    class="delete-btn"
                                                    @click.stop="removeFile(item.name, index, item)"
                                                >
                                                    <el-icon><Close /></el-icon>
                                                </el-button>
                                            </div>
                                            <div class="file-name">{{ f.name }}</div>
                                        </div>
                                        <!-- 继续上传按钮（多文件且已有文件时显示） -->
                                        <div class="continue-upload">
                                            <el-upload
                                                :multiple="true"
                                                :show-file-list="false"
                                                :auto-upload="false"
                                                :on-change="(file, files) => handleFilesChange(files, item.name, item)"
                                            >
                                                <el-button type="primary" size="small">继续上传图片</el-button>
                                            </el-upload>
                                        </div>
                                    </div>
                                    <!-- 单文件预览 -->
                                    <div v-else-if="form[item.name]" class="single-file">
                                        <div class="file-preview-container">
                                            <img v-if="isImageFile({ url: form[item.name] })" :src="form[item.name]" class="file-preview" />
                                            <img v-else src="/image/doc.png" class="file-preview" />
                                            <el-button 
                                                type="danger" 
                                                size="small" 
                                                circle 
                                                class="delete-btn"
                                                @click.stop="removeFile(item.name, 0, item)"
                                            >
                                                <el-icon><Close /></el-icon>
                                            </el-button>
                                        </div>
                                    </div>
                                    <!-- 占位 -->
                                    <div v-else class="upload-placeholder">
                                        <el-icon class="upload-icon">
                                            <Plus />
                                        </el-icon>
                                        <div class="upload-text">点击上传</div>
                                    </div>
                                </el-upload>
                            </template>
                            <template v-else-if="item.type === 'editor'">
                                <div class="editor-container">
                                    <div style="border: 1px solid #ccc;">
                                        <Toolbar 
                                            style="border-bottom: 1px solid #ccc" 
                                            :editor="editorRefs[item.name]" 
                                            :defaultConfig="toolbarConfig" 
                                        />
                                        <Editor
                                            :style="`height: ${parseExtend(item.extend).height || '300px'}; overflow-y: hidden`"
                                            v-model="form[item.name]"
                                            :defaultConfig="editorConfig"
                                            @onCreated="(editor) => handleEditorCreated(editor, item.name)"
                                        />
                                    </div>
                                </div>
                            </template>
                            <template v-else>
                                <el-input v-model="form[item.name]" />
                            </template>
                        </el-form-item>
                    </template>
                    <template v-else>
                        <div>暂无配置项</div>
                    </template>
                </el-col>

                <el-col :span="24">
                    <el-form-item>
                        <el-button type="primary" @click="onSubmit(formRef)">保存</el-button>
                    </el-form-item>
                </el-col>
            </el-row>
        </el-form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts" name="forms">
import service from '@/utils/request';
import { Close, Plus } from '@element-plus/icons-vue';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import '@wangeditor/editor/dist/css/style.css'; // 引入 css
import type { FormInstance, FormProps } from 'element-plus';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import { onBeforeUnmount, reactive, ref } from 'vue';
const labelPosition = ref<FormProps['labelPosition']>()

type ConfigItem = {
    id: number
    name: string
    title: string
    value: string
    group: string
    content: string
    type: string
    extend: string
}

const formRef = ref<FormInstance>();
const groups = ref<string[]>([]);
const currentGroup = ref<string>('');
const items = ref<ConfigItem[]>([]);
const groupedItems = ref<Record<string, ConfigItem[]>>({});
const form = reactive<Record<string, any>>({});

// 存储文件对象的映射
const fileMap = reactive<Record<string, File | File[]>>({});
// 存储文件列表的映射
const fileListMap = reactive<Record<string, Array<{ name: string; url: string; file: File }>>>({});
// 主上传组件 refs，用于清空内部文件列表避免重复上传
const mainUploadRefs = reactive<Record<string, any>>({});
const setMainUploadRef = (prop: string, el: any) => { if (el) mainUploadRefs[prop] = el; };

// 富文本编辑器配置
const editorRefs = reactive<Record<string, any>>({});
const toolbarConfig = {};
const editorConfig = { 
    placeholder: '请输入内容...', 
    readOnly: false
};

// 处理编辑器创建
const handleEditorCreated = (editor: any, prop: string) => {
    editorRefs[prop] = editor;
};

// 组件销毁时，及时销毁编辑器
onBeforeUnmount(() => {
    Object.values(editorRefs).forEach((editor: any) => {
        if (editor && typeof editor.destroy === 'function') {
            editor.destroy();
        }
    });
});

// 解析extend字段
const parseExtend = (extend: string) => {
    try {
        return JSON.parse(extend || '{}');
    } catch (e) {
        return {};
    }
};

// 判断是否为图片文件
const isImageFile = (file: { url: string; name?: string }) => {
    if (!file.url) return false;
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'];
    const url = file.url.toLowerCase();
    const name = (file.name || '').toLowerCase();
    return imageExtensions.some(ext => url.includes(ext) || name.endsWith(ext));
};

// 将值统一为 URL 数组（支持逗号串或数组）
const toUrlArray = (val: any): string[] => {
    if (!val) return [];
    if (Array.isArray(val)) return val.filter(Boolean);
    if (typeof val === 'string') return val.split(',').map(s => s.trim()).filter(s => s.length > 0);
    return [];
};

// 安全 revoke 仅限 blob:
const safeRevoke = (u: string) => { try { if (u?.startsWith('blob:')) URL.revokeObjectURL(u); } catch {} };

// 获取展示文件名
const getBasename = (u: string) => {
    try {
        const qIdx = u.indexOf('?');
        const clean = qIdx >= 0 ? u.slice(0, qIdx) : u;
        return clean.split('/').pop() || '文件';
    } catch { return '文件'; }
};

// 提供给视图的展示列表（优先当次选择，否则使用已保存值）
const getDisplayFileList = (prop: string, item: any): Array<{ name: string; url: string; file?: File }> => {
    const chosen = fileListMap[prop];
    if (chosen && chosen.length) return chosen;
    const urls = toUrlArray((form as any)[prop]);
    return urls.map(u => ({ name: getBasename(u), url: u }));
};

// 删除文件（兼容已保存与当次选择）
const removeFile = (prop: string, index: number, item?: any) => {
    const isMultiple = !!parseExtend(items.value.find(x => x.name === prop)?.extend || '{}').multiple;
    if (isMultiple) {
        if (fileListMap[prop] && fileListMap[prop][index]) {
            safeRevoke(fileListMap[prop][index].url);
            fileListMap[prop].splice(index, 1);
            if (fileListMap[prop].length === 0) {
                delete fileListMap[prop];
                delete fileMap[prop];
                (form as any)[prop] = '';
            } else {
                fileMap[prop] = fileListMap[prop].map(it => it.file as File);
                (form as any)[prop] = fileListMap[prop].map(it => it.url);
            }
        } else {
            const urls = toUrlArray((form as any)[prop]);
            if (index >= 0 && index < urls.length) urls.splice(index, 1);
            (form as any)[prop] = typeof (form as any)[prop] === 'string' ? urls.join(',') : urls;
            if (!urls.length) {
                delete fileMap[prop];
                (form as any)[prop] = '';
            }
        }
    } else {
        if ((form as any)[prop]) {
            safeRevoke((form as any)[prop]);
            delete fileMap[prop];
            (form as any)[prop] = '';
            try { mainUploadRefs[prop]?.clearFiles?.(); } catch {}
        }
    }
};

// 批量上传（一次选择多个合并上传）
const pendingFilesMap: Record<string, any[]> = {};
const uploadTimers: Record<string, any> = {};
const handleFilesChange = (uploadFiles: any[], prop: string, item: any) => {
    pendingFilesMap[prop] = uploadFiles;
    if (uploadTimers[prop]) clearTimeout(uploadTimers[prop]);
    uploadTimers[prop] = setTimeout(async () => {
        const files = (pendingFilesMap[prop] || []).map((uf: any) => uf.raw).filter(Boolean);
        delete pendingFilesMap[prop];
        delete uploadTimers[prop];
        if (!files.length) return;
        try {
            const fd = new FormData();
            files.forEach((f: File, idx: number) => fd.append('files', f, f.name || `file_${idx}`));
            const res = await service.post('/api/api/v1/system/upload', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
            if (res.data?.code === 200) {
                const paths: string[] = res.data?.data?.filePaths || res.data?.data?.paths || [];
                const isMultiple = !!parseExtend(item.extend).multiple;
                if (isMultiple) {
                    const existing = toUrlArray((form as any)[prop]);
                    (form as any)[prop] = [...existing, ...paths];
                    try { mainUploadRefs[prop]?.clearFiles?.(); } catch {}
                } else {
                    (form as any)[prop] = paths[0] || '';
                    try { mainUploadRefs[prop]?.clearFiles?.(); } catch {}
                }
                ElMessage.success('上传成功');
            } else {
                ElMessage.error(res.data?.msg || '上传失败');
            }
        } catch (e: any) {
            ElMessage.error(e?.response?.data?.msg || '上传失败，请重试');
        }
    }, 80);
};

const parseSelectOptions = (content: string): Array<{ label: string; value: string }> => {
    if (!content) return [];
    
    try {
        const parsed = JSON.parse(content);
        if (Array.isArray(parsed)) {
            const options = parsed.map((item) => ({
                label: item,
                value: item
            }));
            return options;
        }
    } catch (e) {
        console.log('Invalid JSON format');
    }
    return [];
}

const initializeFormFromItems = (list: ConfigItem[]) => {
    for (const it of list) {
        if (!(it.name in form)) {
            if (it.type === 'number') {
                const num = Number(it.value)
                form[it.name] = Number.isNaN(num) ? 0 : num
            } else {
                form[it.name] = it.value ?? ''
            }
        } else {
            // 更新已有键的显示值
            if (it.type === 'number') {
                const num = Number(it.value)
                form[it.name] = Number.isNaN(num) ? 0 : num
            } else {
                form[it.name] = it.value ?? ''
            }
        }
    }
}

const getData = async () => {
    NProgress.start();
    try {
        // 首次加载一次即可，返回所有配置项列表
        const res = await service.post('/api/api/v1/config/list', {});
        if (res.data.code === 200) {
            const data = res.data.data || {};
            const list: ConfigItem[] = Array.isArray(data.list) ? data.list : [];

            // 按 group 分组
            const map: Record<string, ConfigItem[]> = {};
            const orderedGroups: string[] = [];
            for (const it of list) {
                const g = it.group || '默认';
                if (!map[g]) {
                    map[g] = [];
                    orderedGroups.push(g);
                }
                map[g].push(it);
            }
            groupedItems.value = map;
            groups.value = orderedGroups;

            // 当前分组与展示项
            if (!currentGroup.value && groups.value.length) {
                currentGroup.value = groups.value[0];
            }
            items.value = groupedItems.value[currentGroup.value] || [];

            // 初始化表单（对所有项初始化，切换 tab 不需再请求）
            initializeFormFromItems(list);
        } else {
            ElMessage.error(res.data.msg || '获取配置失败');
            return false;
        }
    } catch (e) {
        ElMessage.error('未知错误');
        return false;
    } finally {
        NProgress.done();
    }
};

const onGroupChange = () => {
    // 切换分组时，仅在本地切换对应数据，不再请求接口
    items.value = groupedItems.value[currentGroup.value] || [];
}

getData();


// 组合提交数据：将所有 tab 下的数据合并为一个键值对对象，格式 { a:1, b:2 }（包含空值）
const getSubmitData = () => {
    const allItems: ConfigItem[] = Object.values(groupedItems.value || {}).flat();
    const source = allItems.length ? allItems : items.value;
    const payload: Record<string, any> = {};
    source.forEach((item) => {
        let value = (form as any)[item.name];
        // 数组（如多图上传）用逗号拼接为字符串
        if (Array.isArray(value)) {
            value = value.filter(x => x !== null && x !== undefined && x !== '').join(',');
        }
        // 保证包含空值键
        if (value === null || value === undefined) value = '';
        payload[item.name] = value;
    });
    return payload;
};

// 提交：调用配置更新接口（将所有分组的键值对数组一次性提交）
const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    const payload = getSubmitData();
    try {
        // 后端 Update 期望格式: { data: { name: value } }
        const res = await service.post('/api/api/v1/config/update', { data: payload });
        if (res.data?.code === 200) {
            ElMessage.success(res.data?.msg || '保存成功');
            localStorage.setItem('config', JSON.stringify(payload));
            // 通知全局配置已更新（使 header / ucenter 等组件即时刷新）
            window.dispatchEvent(new CustomEvent('configUpdated', { detail: { config: payload } }));
            getData();
        } else {
            ElMessage.error(res.data?.msg || '保存失败');
        }
    } catch (e: any) {
        ElMessage.error(e?.response?.data?.msg || '未知错误');
    }
};
</script>

<style scoped>
.panel {
    margin-bottom: 20px;
    background-color: #fff;
    border: 1px solid transparent;
    border-radius: 4px;
    box-shadow: 0 1px 1px rgba(0, 0, 0, 0.05);
}

.panel-default {
    border-color: #ddd;
}

.panel-heading {
    padding: 10px 15px;
    border-bottom: 1px solid transparent;
    border-top-left-radius: 3px;
    border-top-right-radius: 3px;
    background-color: #f5f5f5;
    border-color: #ddd;
}

.panel-body {
    padding: 15px;
}

.nav-tabs {
    border-bottom: 1px solid #ddd;
    margin-bottom: 0;
    padding-left: 0;
    list-style: none;
}

.nav-tabs > li {
    position: relative;
    display: inline-block;
    margin-bottom: -1px;
}

.nav-tabs > li > a {
    position: relative;
    display: block;
    padding: 10px 15px;
    line-height: 1.42857143;
    color: #555;
    text-decoration: none;
    border: 1px solid transparent;
    border-radius: 4px 4px 0 0;
}

.nav-tabs > li > a:hover {
    border-color: #eee #eee #ddd;
    background-color: #eee;
}

.nav-tabs > li.active > a,
.nav-tabs > li.active > a:hover,
.nav-tabs > li.active > a:focus {
    color: #555;
    cursor: default;
    background-color: #fff;
    border: 1px solid #ddd;
    border-bottom-color: transparent;
}

.nav-tabs > li > a:hover {
    color: #337ab7;
}

/* 上传组件样式 */
.upload-container .el-upload {
    width: 100%;
    border: none;
    background: transparent;
}

.upload-container:hover {
    border-color: var(--el-color-primary);
}

.upload-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
    color: #8c939d;
}

.upload-icon {
    font-size: 28px;
    margin-bottom: 8px;
}

.upload-text {
    font-size: 14px;
}

.single-file {
    display: flex;
    justify-content: center;
    align-items: center;
}

.multiple-files {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
}

.file-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
}

.file-preview-container {
    position: relative;
    display: inline-block;
}

.file-preview {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 4px;
    border: 1px solid #ddd;
}

.delete-btn {
    position: absolute;
    top: -8px;
    right: -8px;
    width: 20px;
    height: 20px;
    min-height: 20px;
    padding: 0;
    font-size: 12px;
    z-index: 10;
}

.delete-btn .el-icon {
    font-size: 12px;
}

.file-name {
    font-size: 12px;
    color: #666;
    margin-top: 4px;
    max-width: 80px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

/* 富文本编辑器样式 */
.editor-container {
    width: 100%;
    margin: 10px 0;
}

.editor-container .w-e-text-container {
    min-height: 200px !important;
}
</style>