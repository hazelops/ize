"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
self["webpackHotUpdate_N_E"]("pages/docs/deploying-ECS-app",{

/***/ "./components/sideBar.js":
/*!*******************************!*\
  !*** ./components/sideBar.js ***!
  \*******************************/
/***/ (function(module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"getStaticProps\": function() { return /* binding */ getStaticProps; },\n/* harmony export */   \"default\": function() { return /* binding */ SideBar; }\n/* harmony export */ });\n/* harmony import */ var C_Users_elect_Desktop_ize_website_node_modules_regenerator_runtime_runtime_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./node_modules/regenerator-runtime/runtime.js */ \"./node_modules/regenerator-runtime/runtime.js\");\n/* harmony import */ var C_Users_elect_Desktop_ize_website_node_modules_regenerator_runtime_runtime_js__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(C_Users_elect_Desktop_ize_website_node_modules_regenerator_runtime_runtime_js__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! react/jsx-runtime */ \"./node_modules/react/jsx-runtime.js\");\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var _ize__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./ize */ \"./components/ize.js\");\n/* harmony import */ var _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../utilities/sideBarMenu */ \"./utilities/sideBarMenu.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! react */ \"./node_modules/react/index.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_4___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_4__);\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! next/link */ \"./node_modules/next/link.js\");\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_5___default = /*#__PURE__*/__webpack_require__.n(next_link__WEBPACK_IMPORTED_MODULE_5__);\n/* module decorator */ module = __webpack_require__.hmd(module);\n\n\n\n\n\n\nfunction asyncGeneratorStep(gen, resolve, reject, _next, _throw, key, arg) {\n    try {\n        var info = gen[key](arg);\n        var value = info.value;\n    } catch (error) {\n        reject(error);\n        return;\n    }\n    if (info.done) {\n        resolve(value);\n    } else {\n        Promise.resolve(value).then(_next, _throw);\n    }\n}\nfunction _asyncToGenerator(fn) {\n    return function() {\n        var self = this, args = arguments;\n        return new Promise(function(resolve, reject) {\n            var gen = fn.apply(self, args);\n            function _next(value) {\n                asyncGeneratorStep(gen, resolve, reject, _next, _throw, \"next\", value);\n            }\n            function _throw(err) {\n                asyncGeneratorStep(gen, resolve, reject, _next, _throw, \"throw\", err);\n            }\n            _next(undefined);\n        });\n    };\n}\nvar _s = $RefreshSig$();\nfunction TopElement(props) {\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"div\", {\n        id: props.title,\n        className: \"flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer\",\n        onClick: props.onClick,\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 8,\n            columnNumber: 9\n        },\n        __self: this,\n        children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"span\", {\n            className: \"mx-4 font-medium capitalize\",\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 9,\n                columnNumber: 13\n            },\n            __self: this,\n            children: props.title\n        })\n    }));\n}\n_c = TopElement;\nfunction NestedMenu(props) {\n    var _this = this;\n    if (props.hidden) {\n        return null;\n    }\n    var nestedList = props.nestedItems.map(function(el) {\n        var pathName = el.slice().replaceAll(\" \", \"-\");\n        var route = pathName == \"welcome\" ? \"\" : pathName;\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)((next_link__WEBPACK_IMPORTED_MODULE_5___default()), {\n            href: \"/docs/\".concat(route),\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 22,\n                columnNumber: 16\n            },\n            __self: _this,\n            children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(TopElement, {\n                title: el,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 23,\n                    columnNumber: 21\n                },\n                __self: _this\n            })\n        }));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"div\", {\n        className: \"flex flex-col justify-between flex-1 ml-10\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 29,\n            columnNumber: 9\n        },\n        __self: this,\n        children: nestedList\n    }));\n}\n_c1 = NestedMenu;\nfunction MenuElement(props) {\n    _s();\n    var ref = (0,react__WEBPACK_IMPORTED_MODULE_4__.useState)(false), isHidden = ref[0], setHidden = ref[1];\n    var handleClick = function handleClick() {\n        setHidden(!isHidden);\n    };\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsxs)((react__WEBPACK_IMPORTED_MODULE_4___default().Fragment), {\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 43,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(TopElement, {\n                title: props.title,\n                onClick: handleClick,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 44,\n                    columnNumber: 13\n                },\n                __self: this\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(NestedMenu, {\n                hidden: isHidden,\n                nestedItems: props.nestedItems,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 45,\n                    columnNumber: 13\n                },\n                __self: this\n            })\n        ]\n    }));\n}\n_s(MenuElement, \"Hdw5EO+DplCNBEJcNuH8tsP7WZ4=\");\n_c2 = MenuElement;\nfunction _getStaticProps() {\n    _getStaticProps = _asyncToGenerator(C_Users_elect_Desktop_ize_website_node_modules_regenerator_runtime_runtime_js__WEBPACK_IMPORTED_MODULE_0___default().mark(function _callee() {\n        var menu;\n        return C_Users_elect_Desktop_ize_website_node_modules_regenerator_runtime_runtime_js__WEBPACK_IMPORTED_MODULE_0___default().wrap(function _callee$(_ctx) {\n            while(1)switch(_ctx.prev = _ctx.next){\n                case 0:\n                    menu = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_3__.mainMenu;\n                    return _ctx.abrupt(\"return\", {\n                        props: {\n                        }\n                    });\n                case 2:\n                case \"end\":\n                    return _ctx.stop();\n            }\n        }, _callee);\n    }));\n    return _getStaticProps.apply(this, arguments);\n}\n//------------------------------------------------------------------------------------------------------------------\nfunction getStaticProps() {\n    return _getStaticProps.apply(this, arguments);\n}\nfunction SideBar() {\n    var _this = this;\n    var menuList = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_3__.mainMenu.map(function(el) {\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(MenuElement, {\n            title: el.title,\n            nestedItems: el.nestedItems,\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 64,\n                columnNumber: 13\n            },\n            __self: _this\n        }));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsxs)(\"div\", {\n        className: \"flex flex-col w-64 h-screen px-4 py-8 bg-white border-r\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 72,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"div\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 73,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(_ize__WEBPACK_IMPORTED_MODULE_2__[\"default\"], {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 74,\n                        columnNumber: 16\n                    },\n                    __self: this\n                })\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"div\", {\n                className: \"flex flex-col justify-between flex-1\",\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 77,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsxs)(\"nav\", {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 78,\n                        columnNumber: 17\n                    },\n                    __self: this,\n                    children: [\n                        menuList,\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(\"hr\", {\n                            className: \"my-6 border-gray-200\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 80,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        }),\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_1__.jsx)(TopElement, {\n                            title: \"See Also\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 81,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        })\n                    ]\n                })\n            })\n        ]\n    }));\n};\n_c3 = SideBar;\nvar _c, _c1, _c2, _c3;\n$RefreshReg$(_c, \"TopElement\");\n$RefreshReg$(_c1, \"NestedMenu\");\n$RefreshReg$(_c2, \"MenuElement\");\n$RefreshReg$(_c3, \"SideBar\");\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,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UyxFQUFDLENBQXlEOzs7Ozs7OztpRkFDbkVILENBQUc7Ozs7Ozs7K0ZBQ0FQLDRDQUFHOzs7Ozs7Ozs7aUZBR05PLENBQUc7Z0JBQUNHLFNBQVMsRUFBQyxDQUFzQzs7Ozs7OztnR0FDaER1QixDQUFHOzs7Ozs7Ozt3QkFDQ0QsUUFBUTs2RkFDUkUsQ0FBRTs0QkFBQ3hCLFNBQVMsRUFBQyxDQUFzQjs7Ozs7Ozs7NkZBQ25DTCxVQUFVOzRCQUFDSSxLQUFLLEVBQUMsQ0FBVTs7Ozs7Ozs7Ozs7OztBQUtoRCxDQUFDO01BMUJ1QnNCLE9BQU8iLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9fTl9FLy4vY29tcG9uZW50cy9zaWRlQmFyLmpzP2ZkZGEiXSwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IEl6ZSBmcm9tIFwiLi9pemVcIlxyXG5pbXBvcnQgeyBtYWluTWVudSB9IGZyb20gJy4uL3V0aWxpdGllcy9zaWRlQmFyTWVudSdcclxuaW1wb3J0IFJlYWN0LCB7IHVzZVN0YXRlIH0gZnJvbSAncmVhY3QnXHJcbmltcG9ydCBMaW5rIGZyb20gJ25leHQvbGluaydcclxuXHJcbmZ1bmN0aW9uIFRvcEVsZW1lbnQocHJvcHMpIHtcclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBpZD17cHJvcHMudGl0bGV9IGNsYXNzTmFtZT1cImZsZXggaXRlbXMtY2VudGVyIHB4LTQgcHktMiBtdC01IHRleHQtZ3JheS02MDAgcm91bmRlZC1tZCBob3ZlcjpiZy1ncmF5LTIwMCB0cmFuc2l0aW9uLWNvbG9ycyBkdXJhdGlvbi0zMDAgdHJhbnNmb3JtIGN1cnNvci1wb2ludGVyXCIgb25DbGljaz17cHJvcHMub25DbGlja30+XHJcbiAgICAgICAgICAgIDxzcGFuIGNsYXNzTmFtZT1cIm14LTQgZm9udC1tZWRpdW0gY2FwaXRhbGl6ZVwiPntwcm9wcy50aXRsZX08L3NwYW4+XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn1cclxuXHJcbmZ1bmN0aW9uIE5lc3RlZE1lbnUocHJvcHMpIHtcclxuICAgIGlmIChwcm9wcy5oaWRkZW4pIHtcclxuICAgICAgICByZXR1cm4gbnVsbFxyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IG5lc3RlZExpc3QgPSBwcm9wcy5uZXN0ZWRJdGVtcy5tYXAoZWwgPT4ge1xyXG4gICAgICAgIGNvbnN0IHBhdGhOYW1lID0gZWwuc2xpY2UoKS5yZXBsYWNlQWxsKFwiIFwiLCBcIi1cIilcclxuICAgICAgICBsZXQgcm91dGUgPSBwYXRoTmFtZSA9PSBcIndlbGNvbWVcIj8gXCJcIiA6IHBhdGhOYW1lXHJcbiAgICAgICAgcmV0dXJuIDxMaW5rIGhyZWY9e2AvZG9jcy8ke3JvdXRlfWB9PlxyXG4gICAgICAgICAgICAgICAgICAgIDxUb3BFbGVtZW50IHRpdGxlPXtlbH0vPlxyXG4gICAgICAgICAgICAgICAgPC9MaW5rPlxyXG4gICAgICAgIFxyXG4gICAgfSlcclxuXHJcbiAgICByZXR1cm4gKFxyXG4gICAgICAgIDxkaXYgY2xhc3NOYW1lPVwiZmxleCBmbGV4LWNvbCBqdXN0aWZ5LWJldHdlZW4gZmxleC0xIG1sLTEwXCI+XHJcbiAgICAgICAgICAgIHtuZXN0ZWRMaXN0fVxyXG4gICAgICAgIDwvZGl2PlxyXG4gICAgKVxyXG59XHJcblxyXG5mdW5jdGlvbiBNZW51RWxlbWVudChwcm9wcykge1xyXG4gICAgY29uc3QgW2lzSGlkZGVuLCBzZXRIaWRkZW5dID0gdXNlU3RhdGUoZmFsc2UpXHJcblxyXG4gICAgY29uc3QgaGFuZGxlQ2xpY2sgPSBmdW5jdGlvbigpIHtcclxuICAgICAgICBzZXRIaWRkZW4oIWlzSGlkZGVuKVxyXG4gICAgfSBcclxuXHJcbiAgICByZXR1cm4gKFxyXG4gICAgICAgIDxSZWFjdC5GcmFnbWVudD5cclxuICAgICAgICAgICAgPFRvcEVsZW1lbnQgdGl0bGU9e3Byb3BzLnRpdGxlfSBvbkNsaWNrPXtoYW5kbGVDbGlja30gLz5cclxuICAgICAgICAgICAgPE5lc3RlZE1lbnUgaGlkZGVuPXtpc0hpZGRlbn0gbmVzdGVkSXRlbXM9e3Byb3BzLm5lc3RlZEl0ZW1zfSAvPlxyXG4gICAgICAgIDwvUmVhY3QuRnJhZ21lbnQ+XHJcbiAgICAgICAgXHJcbiAgICApXHJcbn1cclxuLy8tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS1cclxuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIGdldFN0YXRpY1Byb3BzKCkge1xyXG4gICAgY29uc3QgbWVudSA9IG1haW5NZW51XHJcbiAgICByZXR1cm4ge1xyXG4gICAgICAgIHByb3BzOiB7XHJcblxyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgZnVuY3Rpb24gU2lkZUJhcigpIHtcclxuXHJcbiAgICBjb25zdCBtZW51TGlzdCA9IG1haW5NZW51Lm1hcChlbCA9PiB7XHJcbiAgICAgICAgcmV0dXJuIChcclxuICAgICAgICAgICAgPE1lbnVFbGVtZW50XHJcbiAgICAgICAgICAgICAgICB0aXRsZT17ZWwudGl0bGV9XHJcbiAgICAgICAgICAgICAgICBuZXN0ZWRJdGVtcz17ZWwubmVzdGVkSXRlbXN9XHJcbiAgICAgICAgICAgICAvPlxyXG4gICAgICAgIClcclxuICAgIH0pXHJcblxyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8ZGl2IGNsYXNzTmFtZT1cImZsZXggZmxleC1jb2wgdy02NCBoLXNjcmVlbiBweC00IHB5LTggYmctd2hpdGUgYm9yZGVyLXJcIj5cclxuICAgICAgICAgICAgPGRpdj5cclxuICAgICAgICAgICAgICAgPEl6ZSAvPiBcclxuICAgICAgICAgICAgPC9kaXY+XHJcblxyXG4gICAgICAgICAgICA8ZGl2IGNsYXNzTmFtZT1cImZsZXggZmxleC1jb2wganVzdGlmeS1iZXR3ZWVuIGZsZXgtMVwiPlxyXG4gICAgICAgICAgICAgICAgPG5hdj5cclxuICAgICAgICAgICAgICAgICAgICB7bWVudUxpc3R9XHJcbiAgICAgICAgICAgICAgICAgICAgPGhyIGNsYXNzTmFtZT1cIm15LTYgYm9yZGVyLWdyYXktMjAwXCIgLz5cclxuICAgICAgICAgICAgICAgICAgICA8VG9wRWxlbWVudCB0aXRsZT1cIlNlZSBBbHNvXCIgLz5cclxuICAgICAgICAgICAgICAgIDwvbmF2PlxyXG4gICAgICAgICAgICA8L2Rpdj5cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufSJdLCJuYW1lcyI6WyJJemUiLCJtYWluTWVudSIsIlJlYWN0IiwidXNlU3RhdGUiLCJMaW5rIiwiVG9wRWxlbWVudCIsInByb3BzIiwiZGl2IiwiaWQiLCJ0aXRsZSIsImNsYXNzTmFtZSIsIm9uQ2xpY2siLCJzcGFuIiwiTmVzdGVkTWVudSIsImhpZGRlbiIsIm5lc3RlZExpc3QiLCJuZXN0ZWRJdGVtcyIsIm1hcCIsImVsIiwicGF0aE5hbWUiLCJzbGljZSIsInJlcGxhY2VBbGwiLCJyb3V0ZSIsImhyZWYiLCJNZW51RWxlbWVudCIsImlzSGlkZGVuIiwic2V0SGlkZGVuIiwiaGFuZGxlQ2xpY2siLCJGcmFnbWVudCIsImdldFN0YXRpY1Byb3BzIiwibWVudSIsIlNpZGVCYXIiLCJtZW51TGlzdCIsIm5hdiIsImhyIl0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./components/sideBar.js\n");

/***/ })

});