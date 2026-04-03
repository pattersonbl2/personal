## Summary

Personal website and blog built with Hugo (PaperMod theme), deployed on Cloudflare Pages at [ark31.info](https://ark31.info).

### Tech Stack
- **Frontend:** Hugo static site with PaperMod theme
- **Backend:** Go API on GCP Cloud Run (stage + prod) — handles contact form submissions and resume PDF downloads
- **Infrastructure:** Terraform-managed GCP resources (Cloud Run, Artifact Registry, Secret Manager)
- **Deployment:** Cloudflare Pages — auto-deploys staging on branch commits and production on merge to main

### Changes Made
- **Contact form** — working contact page that submits to the Go backend API with honeypot spam protection
- **Resume page** — full resume rendered in markdown with a PDF download button powered by the backend API
- **Backend API** — Go service deployed to Cloud Run in both stage (`api-stage.ark31.info`) and prod (`api.ark31.info`) environments
- **GCP infrastructure** — Terraform configs for Cloud Run, Artifact Registry, IAM, secrets, and custom domain mappings
- **Profile landing page** — profile mode homepage with image, nav buttons (Archive, Resume, Blog, Contact), and social links (GitHub, LinkedIn)
- **Search and archives** — built-in search and post archive pages
- **Resume PDF sync** — downloadable resume PDF is generated from `content/resume.md` into `backend/resume.pdf` so the website resume and downloaded resume share one source of truth

### Resume PDF Generation
- Local build: `make resume-pdf`
- CI build: `.github/workflows/deploy-backend.yml` regenerates the PDF automatically before backend deploys

### Future Ideas
1. **Blog post series with tutorials** — write technical posts covering real-world SRE/DevOps topics (Terraform patterns, Kubernetes debugging, CI/CD pipelines) to showcase expertise and give back to the community
2. **Visitor analytics dashboard** — integrate a privacy-friendly analytics tool (e.g., Plausible or Umami) to track page views, popular posts, and traffic sources without cookies
3. **Projects showcase page** — a dedicated page highlighting personal and open-source projects with descriptions, tech used, and links to live demos or repos
