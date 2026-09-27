# Gemini 原生生图 API 使用文档

本文档适用于通过 `https://api.zmoapi.cn` 调用 Gemini 原生图片生成协议。

## 1. 接口总览

### Base URL

```text
https://api.zmoapi.cn
```

### 鉴权

所有请求使用中转站分配的 API Key：

```http
Authorization: Bearer YOUR_API_KEY
```

请求 JSON 必须使用 UTF-8 编码：

```http
Content-Type: application/json; charset=utf-8
```

### 接口路径

普通请求：

```http
POST /v1beta/models/{model}:generateContent
```

SSE 流式请求：

```http
POST /v1beta/models/{model}:streamGenerateContent?alt=sse
```

获取当前密钥可用模型：

```http
GET /v1beta/models
```

## 2. 模型名称

推荐使用以下正式模型名：

| 模型 | 定位 | 支持分辨率 |
| --- | --- | --- |
| `gemini-3.1-flash-image` | 速度优先，适合高并发常规生图和图片编辑 | `1K`、`2K`、`4K` |
| `gemini-3-pro-image` | 质量优先，适合复杂提示词和高质量输出 | `1K`、`2K`、`4K` |

以下预览名称作为兼容别名保留：

```text
gemini-3.1-flash-image-preview
gemini-3-pro-image-preview
```

建议客户端先调用模型列表，并使用返回结果中的实际模型名：

```bash
curl 'https://api.zmoapi.cn/v1beta/models' \
  -H 'Authorization: Bearer YOUR_API_KEY'
```

响应示例：

```json
{
  "models": [
    {
      "name": "models/gemini-3.1-flash-image",
      "displayName": "Banana 2",
      "supportedGenerationMethods": [
        "generateContent",
        "streamGenerateContent"
      ]
    }
  ]
}
```

调用模型时，需要去掉返回名称前面的 `models/`：

```text
models/gemini-3.1-flash-image
                  ↓
gemini-3.1-flash-image
```

> `banana2-1k`、`banana2-2k`、`banana2-4k`、`banana-pro-1k`、`banana-pro-2k`、`banana-pro-4k` 属于 OpenAI 图片接口的模型名。调用本文件介绍的 Gemini 原生接口时，应使用上表中的 Gemini 模型名，并通过 `imageSize` 选择分辨率。

## 3. 请求结构

完整请求由两部分组成：

1. `contents`：提示词和参考图。
2. `generationConfig`：返回模态、比例、分辨率和图片交付格式。

基础结构：

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "在这里填写提示词"
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "16:9",
      "imageSize": "2K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

## 4. 分辨率

使用 `generationConfig.imageConfig.imageSize` 控制输出分辨率：

```json
"imageConfig": {
  "imageSize": "2K"
}
```

允许值：

| 参数 | 含义 |
| --- | --- |
| `1K` | 常规分辨率，速度较快 |
| `2K` | 高清分辨率 |
| `4K` | 超高清分辨率，生成时间和资源消耗更高 |

注意事项：

- `K` 必须大写。
- 正确：`"imageSize":"2K"`。
- 错误：`"imageSize":"2k"`。
- 不要传 `2048x2048` 代替 `2K`。
- 实际像素长宽会根据宽高比和模型版本决定，不应假定所有 `2K` 图片都是正方形。

## 5. 宽高比

使用 `generationConfig.imageConfig.aspectRatio`：

```json
"imageConfig": {
  "aspectRatio": "16:9",
  "imageSize": "2K"
}
```

### 通用比例

```text
1:1
16:9
9:16
4:3
3:4
3:2
2:3
5:4
4:5
21:9
```

### Flash 额外支持的超长比例

```text
1:4
4:1
1:8
8:1
```

### 自动比例

希望模型根据参考图自行决定比例时，直接省略 `aspectRatio`：

```json
"imageConfig": {
  "imageSize": "2K"
}
```

不要发送：

```json
"aspectRatio": "auto"
```

`auto` 不是 Gemini `aspectRatio` 的标准枚举值。多张参考图比例不一致时，建议明确指定最终输出比例，不要依靠自动判断。

## 6. 纯文生图

```bash
curl 'https://api.zmoapi.cn/v1beta/models/gemini-3.1-flash-image:generateContent' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -H 'Content-Type: application/json; charset=utf-8' \
  --data-raw '{
    "contents": [
      {
        "role": "user",
        "parts": [
          {
            "text": "生成一张16:9的现代科技产品宣传图，干净明亮，产品位于画面中心"
          }
        ]
      }
    ],
    "generationConfig": {
      "responseModalities": ["TEXT", "IMAGE"],
      "imageConfig": {
        "aspectRatio": "16:9",
        "imageSize": "2K"
      },
      "responseFormat": {
        "image": {
          "delivery": "URI"
        }
      }
    }
  }'
```

## 7. 单张 URL 参考图

公网参考图使用 `fileData`：

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "修改尺寸，保持人物与场景一致性"
        },
        {
          "fileData": {
            "mimeType": "image/png",
            "fileUri": "https://example.com/reference.png"
          }
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "16:9",
      "imageSize": "2K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

URL 必须满足：

- 是公网可访问的 HTTP 或 HTTPS 直链。
- 打开链接后应直接返回图片，而不是网页、登录页或防盗链提示。
- `fileUri` 只能填写纯 URL。
- 不要传 Markdown 链接。

正确：

```json
"fileUri": "https://example.com/reference.png"
```

错误：

```json
"fileUri": "[https://example.com/reference.png](https://example.com/reference.png)"
```

## 8. 单张 Base64 参考图

本地文件需要读取为 Base64，并使用 `inlineData`：

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "保留人物、服装和面部特征，将背景改成海边日落"
        },
        {
          "inlineData": {
            "mimeType": "image/jpeg",
            "data": "BASE64_IMAGE_DATA"
          }
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "16:9",
      "imageSize": "2K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

`inlineData.data` 只能包含纯 Base64 字符串。

正确：

```json
"data": "iVBORw0KGgoAAAANSUhEUgAA..."
```

错误：

```json
"data": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
```

`mimeType` 必须与真实图片类型一致，常用值：

```text
image/png
image/jpeg
image/webp
```

## 9. 多张 URL 参考图

每张参考图对应一个独立的 `fileData` Part：

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "第一张是商品主体，第二张是目标场景，第三张是配色参考。保持第一张商品结构不变，将其自然放入第二张场景，并使用第三张的配色。"
        },
        {
          "fileData": {
            "mimeType": "image/png",
            "fileUri": "https://example.com/product.png"
          }
        },
        {
          "fileData": {
            "mimeType": "image/jpeg",
            "fileUri": "https://example.com/scene.jpg"
          }
        },
        {
          "fileData": {
            "mimeType": "image/webp",
            "fileUri": "https://example.com/color.webp"
          }
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "16:9",
      "imageSize": "4K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

图片在 `parts` 中出现的顺序就是参考图顺序。建议提示词明确说明每张图片的用途，不要只写“参考这些图片”。

## 10. 多张 Base64 参考图

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "第一张作为人物参考，第二张作为服装参考，第三张作为背景参考。"
        },
        {
          "inlineData": {
            "mimeType": "image/png",
            "data": "FIRST_IMAGE_BASE64"
          }
        },
        {
          "inlineData": {
            "mimeType": "image/jpeg",
            "data": "SECOND_IMAGE_BASE64"
          }
        },
        {
          "inlineData": {
            "mimeType": "image/webp",
            "data": "THIRD_IMAGE_BASE64"
          }
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "9:16",
      "imageSize": "2K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

## 11. URL 与 Base64 混合参考图

同一次请求可以混合使用 `fileData` 和 `inlineData`：

```json
{
  "contents": [
    {
      "role": "user",
      "parts": [
        {
          "text": "第一张为人物，第二张为本地服装参考，第三张为场景参考。"
        },
        {
          "fileData": {
            "mimeType": "image/png",
            "fileUri": "https://example.com/person.png"
          }
        },
        {
          "inlineData": {
            "mimeType": "image/jpeg",
            "data": "CLOTHING_IMAGE_BASE64"
          }
        },
        {
          "fileData": {
            "mimeType": "image/webp",
            "fileUri": "https://example.com/scene.webp"
          }
        }
      ]
    }
  ],
  "generationConfig": {
    "responseModalities": ["TEXT", "IMAGE"],
    "imageConfig": {
      "aspectRatio": "9:16",
      "imageSize": "4K"
    },
    "responseFormat": {
      "image": {
        "delivery": "URI"
      }
    }
  }
}
```

一个 Part 只能包含一种数据类型。不要在同一个 Part 中同时放 `fileData` 和 `inlineData`。

## 12. 图片返回方式

### 返回 URL（推荐）

```json
"responseFormat": {
  "image": {
    "delivery": "URI"
  }
}
```

生成结果位于 `fileData.fileUri`：

```json
{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [
          {
            "text": "图片已生成。"
          },
          {
            "fileData": {
              "mimeType": "image/jpeg",
              "fileUri": "https://cdn.example.com/generated-image.jpg"
            }
          }
        ]
      },
      "finishReason": "STOP"
    }
  ],
  "usageMetadata": {
    "promptTokenCount": 100,
    "candidatesTokenCount": 200,
    "totalTokenCount": 300
  }
}
```

### 返回 Base64

确实需要 Base64 时可以指定：

```json
"responseFormat": {
  "image": {
    "delivery": "INLINE"
  }
}
```

生成结果位于 `inlineData.data`：

```json
{
  "inlineData": {
    "mimeType": "image/jpeg",
    "data": "BASE64_GENERATED_IMAGE"
  }
}
```

URL 和 Base64 二选一。一般业务推荐 `URI`，可以减少响应体积和内存占用。

## 13. JavaScript 完整调用示例

以下函数支持纯文生图、URL 参考图、本地参考图以及两者混合。

```javascript
async function fileToInlinePart(file) {
  const dataUrl = await new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result));
    reader.onerror = () => reject(reader.error || new Error("读取图片失败"));
    reader.readAsDataURL(file);
  });

  const comma = dataUrl.indexOf(",");
  if (comma < 0) throw new Error("无法解析图片 Base64");

  return {
    inlineData: {
      mimeType: file.type || "image/png",
      data: dataUrl.slice(comma + 1)
    }
  };
}

function urlToFilePart(url, mimeType = "image/png") {
  return {
    fileData: {
      mimeType,
      fileUri: url
    }
  };
}

async function generateGeminiImage({
  apiKey,
  model = "gemini-3.1-flash-image",
  prompt,
  references = [],
  aspectRatio = "16:9",
  imageSize = "2K",
  delivery = "URI"
}) {
  const referenceParts = await Promise.all(
    references.map(async (reference) => {
      if (reference instanceof File || reference instanceof Blob) {
        return fileToInlinePart(reference);
      }

      if (typeof reference === "string") {
        return urlToFilePart(reference);
      }

      if (reference && typeof reference.url === "string") {
        return urlToFilePart(reference.url, reference.mimeType);
      }

      throw new Error("不支持的参考图类型");
    })
  );

  const imageConfig = { imageSize };
  if (aspectRatio && aspectRatio !== "auto") {
    imageConfig.aspectRatio = aspectRatio;
  }

  const response = await fetch(
    `https://api.zmoapi.cn/v1beta/models/${encodeURIComponent(model)}:generateContent`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${apiKey}`,
        "Content-Type": "application/json; charset=utf-8"
      },
      body: JSON.stringify({
        contents: [
          {
            role: "user",
            parts: [
              { text: prompt },
              ...referenceParts
            ]
          }
        ],
        generationConfig: {
          responseModalities: ["TEXT", "IMAGE"],
          imageConfig,
          responseFormat: {
            image: { delivery }
          }
        }
      })
    }
  );

  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(payload.error?.message || `HTTP ${response.status}`);
  }

  const images = [];
  const texts = [];

  for (const candidate of payload.candidates || []) {
    for (const part of candidate.content?.parts || []) {
      // thought 是模型内部思考内容，不应当作最终结果展示。
      if (part.thought === true) continue;

      if (typeof part.text === "string") {
        texts.push(part.text);
      }

      if (
        part.fileData?.fileUri &&
        part.fileData?.mimeType?.startsWith("image/")
      ) {
        images.push({
          type: "url",
          mimeType: part.fileData.mimeType,
          url: part.fileData.fileUri
        });
      }

      if (
        part.inlineData?.data &&
        part.inlineData?.mimeType?.startsWith("image/")
      ) {
        images.push({
          type: "base64",
          mimeType: part.inlineData.mimeType,
          data: part.inlineData.data,
          dataUrl: `data:${part.inlineData.mimeType};base64,${part.inlineData.data}`
        });
      }
    }
  }

  return { images, texts, raw: payload };
}
```

单张本地参考图：

```javascript
const result = await generateGeminiImage({
  apiKey: "YOUR_API_KEY",
  model: "gemini-3.1-flash-image",
  prompt: "保持人物一致，将画面改成16:9电影场景",
  references: [document.querySelector("#file").files[0]],
  aspectRatio: "16:9",
  imageSize: "2K"
});

console.log(result.images);
```

多张参考图：

```javascript
const result = await generateGeminiImage({
  apiKey: "YOUR_API_KEY",
  model: "gemini-3-pro-image",
  prompt: "第一张作为商品主体，第二张作为背景，第三张作为配色参考。",
  references: [
    document.querySelector("#product").files[0],
    "https://example.com/background.jpg",
    document.querySelector("#colors").files[0]
  ],
  aspectRatio: "16:9",
  imageSize: "4K"
});

console.log(result.images);
```

## 14. SSE 流式调用

流式接口：

```text
POST https://api.zmoapi.cn/v1beta/models/{model}:streamGenerateContent?alt=sse
```

请求体与普通 `generateContent` 完全相同，请求头建议增加：

```http
Accept: text/event-stream
```

示例：

```bash
curl -N 'https://api.zmoapi.cn/v1beta/models/gemini-3.1-flash-image:streamGenerateContent?alt=sse' \
  -H 'Authorization: Bearer YOUR_API_KEY' \
  -H 'Content-Type: application/json; charset=utf-8' \
  -H 'Accept: text/event-stream' \
  --data-raw '{
    "contents": [
      {
        "role": "user",
        "parts": [
          {"text":"生成一张16:9的现代产品宣传图"}
        ]
      }
    ],
    "generationConfig": {
      "responseModalities": ["TEXT", "IMAGE"],
      "imageConfig": {
        "aspectRatio": "16:9",
        "imageSize": "2K"
      },
      "responseFormat": {
        "image": {
          "delivery": "URI"
        }
      }
    }
  }'
```

服务会返回多个 SSE 事件，每个 `data:` 后面是一段 Gemini JSON。客户端应完整解析 SSE 事件，不要假定每一行就是完整 JSON。

## 15. 错误响应

错误遵循 Gemini 风格：

```json
{
  "error": {
    "code": 400,
    "message": "具体错误信息",
    "status": "INVALID_ARGUMENT"
  }
}
```

常见状态：

| HTTP | 状态 | 常见原因 |
| --- | --- | --- |
| `400` | `INVALID_ARGUMENT` | JSON 错误、字段错误、比例不支持、Base64 无效、非 UTF-8 请求 |
| `401` | `UNAUTHENTICATED` | API Key 缺失或无效 |
| `403` | `PERMISSION_DENIED` | 密钥无权限访问模型 |
| `404` | `NOT_FOUND` | 模型名或接口路径错误 |
| `429` | `RESOURCE_EXHAUSTED` | 请求过快、当前额度或并发不足 |
| `500` | `INTERNAL` | 服务内部错误 |
| `503` | `UNAVAILABLE` | 暂无可用渠道或上游临时不可用 |

调用失败时应展示 `error.message`，不要只展示 HTTP 状态码。

## 16. 重要规则

1. Gemini 原生协议没有 OpenAI 图片接口的 `n` 参数。
2. 一次调用通常生成一张最终图片；需要多次生成时，由客户端分别发送多个请求。
3. 参考图必须放在 `contents[].parts[]` 中，不能使用 `imageUrls`。
4. URL 图片使用 `fileData.fileUri`。
5. 本地图片使用 `inlineData.data`，内容必须是纯 Base64。
6. 多图顺序以 `parts` 中的出现顺序为准。
7. 每个 Part 只能包含一种数据类型。
8. 自动比例必须省略 `aspectRatio`，不能传 `auto`。
9. `imageSize` 中的 `K` 必须大写。
10. 最终图片只读取非 `thought` Part 中的 `fileData` 或 `inlineData`。
11. 不要把响应文本中的普通 URL 当作生成图片。
12. 不要因为超时或连接断开自动重试生图请求，以免产生额外生成任务。

## 17. 官方参考

- [Gemini 图片生成文档](https://ai.google.dev/gemini-api/docs/generate-content/image-generation)
- [Gemini GenerateContent API 参考](https://ai.google.dev/api/generate-content)

