export const contentToUrl = (baseUrl: string, value: string) =>
  `${baseUrl}?content=${encodeURIComponent(value)}`;

export const urlToContent = (url: string) => {
  const contentStr = '?content=';

  if (url.indexOf(contentStr) > -1) {
    return decodeURIComponent(
      url.substring(url.indexOf(contentStr) + contentStr.length)
    );
  } else {
    return null;
  }
};
