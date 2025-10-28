<template>
	<el-form ref="formRef" :model="form" :rules="rules" :label-width="options.labelWidth">
		<el-row>
			<el-col :span="item.span || options.span" v-for="item in options.list" :key="item.prop">
				<el-form-item :label="item.label" :prop="item.prop">
					<!-- 文本框、数字框、下拉框、日期框、开关、上传 -->
					<el-input v-if="item.type === 'input'" v-model="form[item.prop]" :disabled="item.disabled || !update" 
						:placeholder="item.placeholder" clearable :type="item.inputType === 'password' ? 'password' : 'text'">
						<template #suffix v-if="item.suffixLink">
						<el-tooltip :content="item.suffixTooltip || '点击跳转'" placement="top">
						<el-icon
							style="cursor: pointer"
							@click.stop="handleSuffixClick(item.suffixLink)"
						>
							<Link />
						</el-icon>
						</el-tooltip>
					</template>
					
					</el-input>
					<el-input-number v-else-if="item.type === 'number'" v-model="form[item.prop]"
						:disabled="item.disabled || !update" controls-position="right"></el-input-number>
					<el-select v-else-if="item.type === 'select'" v-model="form[item.prop]" :disabled="item.disabled || !update"
						:placeholder="item.placeholder" clearable>
						<el-option v-for="opt in item.opts" :label="opt.label" :value="opt.value"></el-option>
					</el-select>
					<el-date-picker v-else-if="item.type === 'date'" type="date" v-model="form[item.prop]"
					:value-format="item.format" :disabled="!update"></el-date-picker>
					<el-date-picker v-else-if="item.type === 'daterange'" type="daterange" v-model="form[item.prop]"
						:value-format="item.format" range-separator="至" start-placeholder="开始日期" 
						end-placeholder="结束日期" :disabled="item.disabled || !update" clearable></el-date-picker>
					<el-switch v-else-if="item.type === 'switch'" v-model="form[item.prop]"
						:active-value="item.activeValue" :inactive-value="item.inactiveValue"
						:active-text="item.activeText" :inactive-text="item.inactiveText"
						:disabled="!update">
					</el-switch>
					<!-- 上传（支持单/多文件，非图片显示 doc 图标） -->
					<el-upload v-else-if="item.type === 'upload'" 
						class="upload-container" 
						style="border: 1px dashed #ccc; padding: 10px; border-radius: 6px;"
						:auto-upload="false"
						:show-file-list="false"
						:multiple="item.multiple || false"
						:disabled="!update || (item.multiple && getDisplayFileList(item.prop, item).length > 0)"
						:on-change="(file, files) => handleFilesChange(files, item.prop, item)"
						:ref="el => setMainUploadRef(item.prop, el)"
					>
						<!-- 多文件预览 -->
						<div v-if="item.multiple && getDisplayFileList(item.prop, item).length" class="multiple-files">
							<div v-for="(f, idx) in getDisplayFileList(item.prop, item)" :key="idx" class="file-item">
								<div class="file-preview-container">
									<img v-if="isImageFile(f)" :src="f.url" class="file-preview" />
									<img v-else src="/image/doc.png" class="file-preview" />
									<el-button 
										type="danger" 
										size="small" 
										circle 
										class="delete-btn"
										@click.stop="removeFile(item.prop, idx, item)"
										v-if="update"
									>
										<el-icon><Close /></el-icon>
									</el-button>
								</div>
								<div class="file-name">{{ f.name }}</div>
							</div>
							<!-- 继续上传按钮（仅在多文件且已有文件时显示） -->
							<div class="continue-upload" v-if="update">
								<el-upload
									:multiple="true"
									:show-file-list="false"
									:auto-upload="false"
									:disabled="!update"
									:on-change="(file, files) => handleFilesChange(files, item.prop, item)"
									
								>
									<el-button type="primary" size="small">继续上传</el-button>
								</el-upload>
							</div>
						</div>
						<!-- 单文件预览 -->
						<div v-else-if="form[item.prop]" class="single-file">
							<div class="file-preview-container">
								<img v-if="isSingleImage(item.prop) || isImageFile({ url: form[item.prop] })" :src="form[item.prop]" class="file-preview" />
								<img v-else src="/image/doc.png" class="file-preview" />
								<el-button 
									type="danger" 
									size="small" 
									circle 
									class="delete-btn"
									@click.stop="removeFile(item.prop, 0, item)"
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
					<!-- 富文本编辑器 -->
					<div v-else-if="item.type === 'editor'" class="editor-container">
						<div style="border: 1px solid #ccc;">
							<Toolbar 
								style="border-bottom: 1px solid #ccc" 
								:editor="editorRefs[item.prop]" 
								:defaultConfig="toolbarConfig" 
							/>
							<Editor
								:style="`height: ${item.height || '300px'}; overflow-y: hidden`"
								v-model="form[item.prop]"
								:defaultConfig="editorConfig"
								@onCreated="(editor) => handleEditorCreated(editor, item.prop)"
							/>
						</div>
					</div>
					<slot :name="item.prop" v-else>

					</slot>
				</el-form-item>
			</el-col>
		</el-row>

		<el-form-item v-if="update">
			<el-button type="primary" @click="saveEdit(formRef)">保 存</el-button>
		</el-form-item>
	</el-form>
</template>

<script lang="ts" setup>
import { FormOption } from '@/types/form-option';
import service from '@/utils/request';
import { Close, Plus } from '@element-plus/icons-vue';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import '@wangeditor/editor/dist/css/style.css'; // 引入 css
import { ElMessage, FormInstance, FormRules } from 'element-plus';
import { onBeforeUnmount, PropType, reactive, ref, watch } from 'vue';
const handleSuffixClick = (url: string) => {
  if (typeof window !== 'undefined') {
    window.open(url, '_blank');
  }
}
const { options, formData, edit, update } = defineProps({
	options: {
		type: Object as PropType<FormOption>,
		required: true
	},
	formData: {
		type: Object,
		required: true
	},
	edit: {
		type: Boolean,
		required: false
	},
	update: {
		type: Function,
		required: false
	}
});


const form = ref({ ...formData });

// 富文本编辑器配置
const editorRefs = reactive<Record<string, any>>({});
const toolbarConfig = {};
const editorConfig = { 
	placeholder: '请输入内容...', 
	readOnly: !update // 根据update参数控制只读状态
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

// 监听 formData 变化，确保每次弹窗打开都能拿到最新默认值
watch(() => formData, (val) => {
  form.value = { ...val };
}, { immediate: true });

const rules: FormRules = options.list.map(item => {
	if (item.required) {
		return { [item.prop]: [{ required: true, message: `${item.label}不能为空`, trigger: 'blur' }] };
	}
	return {};
}).reduce((acc, cur) => ({ ...acc, ...cur }), {});


const formRef = ref<FormInstance>();
const saveEdit = (formEl: FormInstance | undefined) => {
	if (!formEl || !update) return;
	formEl.validate(valid => {
		if (!valid) return false;
		// 使用getSubmitData获取包含文件对象的数据
		update(getSubmitData());
	});
};

// ============== 上传相关 ==============
// 存储文件对象（单或多）的映射
const fileMap = reactive<Record<string, File | File[]>>({});
// 多文件时用于预览的文件列表
const fileListMap = reactive<Record<string, Array<{ name: string; url: string; file: File }>>>({});
// 主上传组件 refs（按字段存储），用于在单文件上传成功后清空内部文件列表，避免重复上传
const mainUploadRefs = reactive<Record<string, any>>({});
const setMainUploadRef = (prop: string, el: any) => { if (el) mainUploadRefs[prop] = el; };

// 是否为图片文件（根据后缀简单判断）
const isImageFile = (file: { url: string; name?: string }) => {
	if (!file.url) return false;
	const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'];
	const url = file.url.toLowerCase();
	const name = (file.name || '').toLowerCase();
	return imageExtensions.some(ext => url.includes(ext) || name.endsWith(ext));
};

// 单文件场景：优先用真实文件 MIME 判断（blob: URL 无后缀会导致扩展名判断失败）
const isSingleImage = (prop: string) => {
	const f = fileMap[prop];
	if (f && f instanceof File) {
		return (f as File).type?.startsWith('image/');
	}
	return false;
};

// 获取文件列表（多文件，原始：仅返回选择的文件列表）
const getFileList = (prop: string) => {
	return fileListMap[prop] || [];
};

// 将字符串或数组值统一转为 URL 数组
const toUrlArray = (val: any): string[] => {
	if (!val) return [];
	if (Array.isArray(val)) return val.filter(Boolean);
	if (typeof val === 'string') {
		return val.split(',').map(s => s.trim()).filter(s => s.length > 0);
	}
	return [];
};

// 提供给视图层的多文件展示列表（兼容“已保存的逗号字符串/数组”与“本次新选文件”）
const getDisplayFileList = (prop: string, item: any): Array<{ name: string; url: string; file?: File }> => {
	const chosen = fileListMap[prop];
	if (chosen && chosen.length) return chosen;
	// 初始值来源于表单：可能是逗号分隔的字符串或数组
	const urls = toUrlArray((form.value as any)[prop]);
	return urls.map(u => ({ name: getBasename(u), url: u }));
};

// 安全释放仅限 blob: 的 URL
const safeRevoke = (u: string) => {
	try { if (u && u.startsWith('blob:')) URL.revokeObjectURL(u); } catch {}
};

// 获取 URL 基名作为文件名展示
const getBasename = (u: string) => {
	try {
		const qIdx = u.indexOf('?');
		const clean = qIdx >= 0 ? u.slice(0, qIdx) : u;
		return clean.split('/').pop() || '文件';
	} catch {
		return '文件';
	}
};

// 删除文件
const removeFile = (prop: string, index: number, item: any) => {
	const isMultiple = !!item.multiple;
	if (isMultiple) {
		if (fileListMap[prop] && fileListMap[prop][index]) {
			// 来自当次选择的文件
			safeRevoke(fileListMap[prop][index].url);
			fileListMap[prop].splice(index, 1);
			if (fileListMap[prop].length === 0) {
				delete fileListMap[prop];
				delete fileMap[prop];
				(form.value as any)[prop] = '';
			} else {
				fileMap[prop] = fileListMap[prop].map(it => it.file as File);
				(form.value as any)[prop] = fileListMap[prop].map(it => it.url);
			}
		} else {
			// 初始值来自已保存的 URL（字符串/数组）
			const urls = toUrlArray((form.value as any)[prop]);
			if (index >= 0 && index < urls.length) {
				urls.splice(index, 1);
				(form.value as any)[prop] = typeof (form.value as any)[prop] === 'string' ? urls.join(',') : urls;
			}
			if (urls.length === 0) {
				delete fileMap[prop];
				(form.value as any)[prop] = '';
			}
		}
	} else {
		if (fileMap[prop]) {
			if (typeof (form.value as any)[prop] === 'string') {
				safeRevoke((form.value as any)[prop]);
			}
			delete fileMap[prop];
			(form.value as any)[prop] = '';
		} else {
			// 初始值为已保存的 URL
			(form.value as any)[prop] = '';
		}
	}
};

// 批量上传节流：同一属性的选择在短时间内合并为一次上传
const pendingFilesMap: Record<string, any[]> = {};
const uploadTimers: Record<string, any> = {};

// 选择文件后（支持一次选择多个），合并为一次请求上传
const handleFilesChange = (uploadFiles: any[], prop: string, item: any) => {
	pendingFilesMap[prop] = uploadFiles;
	if (uploadTimers[prop]) {
		clearTimeout(uploadTimers[prop]);
	}
	uploadTimers[prop] = setTimeout(async () => {
		const files = (pendingFilesMap[prop] || []).map((uf: any) => uf.raw).filter(Boolean);
		delete pendingFilesMap[prop];
		delete uploadTimers[prop];
		if (!files.length) return;
		try {
			const formData = new FormData();
			files.forEach((f: File, idx: number) => formData.append('files', f, f.name || `file_${idx}`));
			const res = await service.post('/api/api/v1/system/upload', formData, {
				headers: { 'Content-Type': 'multipart/form-data' }
			});
			if (res.data?.code === 200) {
				const paths: string[] = res.data?.data?.filePaths || res.data?.data?.paths || [];
				const isMultiple = !!item.multiple;
				if (isMultiple) {
					const existing = toUrlArray((form.value as any)[prop]);
					(form.value as any)[prop] = [...existing, ...paths];
					// 清空（继续上传按钮实例会是独立的el-upload，但也建议清空其内部列表避免重复）
					try { mainUploadRefs[prop]?.clearFiles?.(); } catch {}
				} else {
					(form.value as any)[prop] = paths[0] || '';
					// 单文件：清空主上传组件内部文件记录，避免再次点击时重复上传之前的文件
					try { mainUploadRefs[prop]?.clearFiles?.(); } catch {}
				}
				ElMessage.success('上传成功');
			} else {
				ElMessage.error(res.data?.msg || '上传失败');
			}
		} catch (e: any) {
			ElMessage.error(e?.response?.data?.msg || '上传失败，请重试');
		}
	}, 100);
};

// 获取实际要提交的表单数据（将预览URL替换为真实文件或文件数组）
const getSubmitData = () => {
	const submitData: Record<string, any> = { ...(form.value as any) };
	Object.keys(fileMap).forEach(prop => {
		if (fileMap[prop]) {
			submitData[prop] = fileMap[prop];
		}
	});
	return submitData;
};
// ============== 上传相关 END ==============

</script>

<style>
/* 复用并扩展上传样式，支持单/多文件预览 */
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

/* 富文本编辑器样式（保留） */
.avatar-uploader .el-upload {
	border: 1px dashed var(--el-border-color);
	border-radius: 6px;
	cursor: pointer;
	position: relative;
	overflow: hidden;
	transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
	border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
	font-size: 28px;
	color: #8c939d;
	width: 178px;
	height: 178px;
	text-align: center;
}

.avatar {
	width: 178px;
	height: 178px;
	display: block;
	object-fit: cover;
}

.editor-container {
	width: 100%;
	margin: 10px 0;
}

.editor-container .w-e-text-container {
	min-height: 200px !important;
}

/* 禁用状态样式 */
.avatar-uploader.is-disabled {
	cursor: not-allowed;
	opacity: 0.6;
}

.avatar-uploader.is-disabled .el-upload {
	cursor: not-allowed;
	pointer-events: none;
}
</style>
