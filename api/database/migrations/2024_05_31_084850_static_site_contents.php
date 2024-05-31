<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('static_site_contents', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->foreignId('lang_code')->constrained(table: 'languages', column: 'code')->onUpdate('cascade')->onDelete('cascade');
            $table->string('content');
            $table->unique(['lang_code', 'content']);

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('static_site_contents');
    }
};
