import{S as L,i as Q,s as X,k as c,q as N,a as U,l as i,m,r as T,h as n,c as I,n as t,b as j,J as a,K as Y,W as Z,I as W}from"../../../chunks/index-a598add7.js";import{g as $}from"../../../chunks/navigation-cb28aaf0.js";import{f as z}from"../../../chunks/FeedbackStore-f117c57a.js";function ee(y){let u,h,l,p,e,s,_,G,A,o,g,D,v,E,V,q,d,k,C,b,w,O,R,H;return{c(){u=c("h1"),h=N("Goly -- Update"),l=U(),p=c("div"),e=c("form"),s=c("div"),_=c("span"),G=N("Redirect to"),A=U(),o=c("input"),D=U(),v=c("div"),E=c("span"),V=N("Goly"),q=U(),d=c("input"),C=U(),b=c("div"),w=c("button"),O=N("Update"),this.h()},l(r){u=i(r,"H1",{class:!0});var f=m(u);h=T(f,"Goly -- Update"),f.forEach(n),l=I(r),p=i(r,"DIV",{class:!0});var J=m(p);e=i(J,"FORM",{class:!0});var x=m(e);s=i(x,"DIV",{class:!0});var P=m(s);_=i(P,"SPAN",{});var B=m(_);G=T(B,"Redirect to"),B.forEach(n),A=I(P),o=i(P,"INPUT",{type:!0,class:!0,placeholder:!0,name:!0,autocomplete:!0}),P.forEach(n),D=I(x),v=i(x,"DIV",{class:!0});var S=m(v);E=i(S,"SPAN",{});var F=m(E);V=T(F,"Goly"),F.forEach(n),q=I(S),d=i(S,"INPUT",{type:!0,class:!0,placeholder:!0,name:!0,autocomplete:!0}),S.forEach(n),C=I(x),b=i(x,"DIV",{class:!0});var K=m(b);w=i(K,"BUTTON",{class:!0});var M=m(w);O=T(M,"Update"),M.forEach(n),K.forEach(n),x.forEach(n),J.forEach(n),this.h()},h(){t(u,"class","text-3xl text-sky-500 my-5 text-center"),t(o,"type","text"),t(o,"class","border border-sky-500 rounded-md p-1 w-full"),t(o,"placeholder","https://www.bbc.co.uk"),o.value=g=y[0].redirect,t(o,"name","redirect"),o.required=!0,t(o,"autocomplete","off"),t(s,"class","flex flex-col w-full py-2"),t(d,"type","text"),t(d,"class","border border-sky-500 rounded-md p-1 w-full"),t(d,"placeholder","Short link or leave blank to have a random one generated"),d.value=k=y[0].goly,t(d,"name","goly"),t(d,"autocomplete","off"),t(v,"class","flex flex-col w-full py-2"),t(w,"class","text-white bg-sky-500 rounded-md px-3 py-2 w-full"),t(b,"class","py-5"),t(e,"class","min-w-full"),t(p,"class","flex flex-col mx-auto w-5/6 md:w-1/2 lg:w-1/3 border border-slate-500 rounded-md p-2")},m(r,f){j(r,u,f),a(u,h),j(r,l,f),j(r,p,f),a(p,e),a(e,s),a(s,_),a(_,G),a(s,A),a(s,o),a(e,D),a(e,v),a(v,E),a(E,V),a(v,q),a(v,d),a(e,C),a(e,b),a(b,w),a(w,O),R||(H=Y(e,"submit",Z(y[1])),R=!0)},p(r,[f]){f&1&&g!==(g=r[0].redirect)&&o.value!==g&&(o.value=g),f&1&&k!==(k=r[0].goly)&&d.value!==k&&(d.value=k)},i:W,o:W,d(r){r&&n(u),r&&n(l),r&&n(p),R=!1,H()}}}function te(y,u,h){let{data:l}=u;async function p(e){try{h(0,l.redirect=e.target[0].value,l),h(0,l.goly=e.target[1].value,l),h(0,l.random=e.target[1].value==="",l);const s=await fetch("/api/goly",{method:"PATCH",body:JSON.stringify(l),headers:{"Content-Type":"application/json"}});if(!s.ok)throw new Error;const _=await s.json();z.set("Success: Goly has been successfully updated."),$("/")}catch{z.set("Error: An error has occurred. Please try again.")}}return y.$$set=e=>{"data"in e&&h(0,l=e.data)},[l,p]}class se extends L{constructor(u){super(),Q(this,u,te,ee,X,{data:0})}}export{se as default};
